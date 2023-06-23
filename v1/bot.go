package telebot

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultApiUrl       string        = "https://api.telebot.org"
	DefaultSendTimeout  time.Duration = 2 * time.Second
	DefultUpdateTimeout time.Duration = 2 * time.Second
)

type ErrStatus struct {
	ErrorCode   int
	Description string
}

func (se ErrStatus) Error() string {
	return fmt.Sprintf("telegram status error '%s', error_code: %d", se.Description, se.ErrorCode)
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Request interface {
	GetParams() (v url.Values, method string, err error)
}

type Response interface {
	Parse(reader io.Reader) error
}

type UpdatesRequest struct {
	Offset         int      `json:"offset"`
	Limit          int      `json:"limit"`
	Timeout        int      `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

func (req UpdatesRequest) GetParams() (val url.Values, method string, err error) {
	method = "getUpdates"
	val = url.Values{}
	if req.Offset != 0 {
		val.Add("offset", strconv.Itoa(req.Offset))
	}
	if req.Limit > 0 {
		val.Add("limit", strconv.Itoa(req.Limit))
	}
	if req.Timeout > 0 {
		val.Add("timeout", strconv.Itoa(req.Timeout))
	}
	for _, au := range req.AllowedUpdates {
		val.Add("allowed_updates", au)
	}
	return
}

func NewSimpleBot(Token string, client httpClient) SimpleBot {
	return SimpleBot{
		apiEndpoint:   DefaultApiUrl,
		client:        client,
		token:         Token,
		sendTimeout:   DefaultSendTimeout,
		updateTimeout: DefultUpdateTimeout,
	}
}

type SimpleBot struct {
	apiEndpoint   string
	client        httpClient
	token         string
	sendTimeout   time.Duration
	updateTimeout time.Duration
}

func (sb SimpleBot) sendRequest(ctx context.Context, req Request) (*http.Response, error) {
	values, method, err := req.GetParams()

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/bot%s/%s", sb.apiEndpoint, sb.token, method)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(values.Encode()))

	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpResp, err := sb.client.Do(httpReq)

	return httpResp, err
}

func (sb SimpleBot) GetUpdates(ctx context.Context, req UpdatesRequest) (UpdateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, sb.updateTimeout)
	defer cancel()
	httpResp, err := sb.sendRequest(ctx, req)
	if err != nil {
		return UpdateResponse{}, err
	}
	defer httpResp.Body.Close()

	ur := UpdateResponse{}

	if err = ur.Parse(httpResp.Body); err != nil {
		return UpdateResponse{}, err
	}

	if !ur.Ok {
		return UpdateResponse{}, ErrStatus{ErrorCode: ur.ErrorCode, Description: ur.Description}
	}

	return ur, nil
}

func (sb SimpleBot) Send(ctx context.Context, req Request) (MessageResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, sb.sendTimeout)
	defer cancel()

	httpResp, err := sb.sendRequest(ctx, req)
	if err != nil {
		return MessageResponse{}, err
	}
	defer httpResp.Body.Close()
	mr := MessageResponse{}
	err = mr.Parse(httpResp.Body)
	return mr, err
}

type SimplePoller struct {
	bot           Bot
	offset        int
	updateHandler UpdateHandler
}

func (slp SimplePoller) getUpdates(ctx context.Context) (UpdateResponse, error) {
	resp, err := slp.bot.GetUpdates(ctx, UpdatesRequest{Offset: slp.offset})
	if err != nil {
		return UpdateResponse{}, fmt.Errorf("get updates with offset %d error: '%w'", slp.offset, err)
	}

	return resp, nil
}

func (sp SimplePoller) ProceedUpdates(ctx context.Context) (offset int, err error) {
	ur, err := sp.getUpdates(ctx)

	if err != nil {
		return sp.offset, err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, update := range ur.Result {
		err := sp.updateHandler.Proceed(ctx, sp.bot, update)
		if err != nil {
			return sp.offset, fmt.Errorf("proceed update %v error: '%w'", update, err)
		}
		sp.offset = update.UpdateId + 1
	}
	return sp.offset, nil
}

func NewSimplePoller(b Bot, handler UpdateHandler) SimplePoller {
	return SimplePoller{bot: b, updateHandler: handler}
}
