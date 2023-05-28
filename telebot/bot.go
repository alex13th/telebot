package telebot

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

const DefaultApiUrl string = "https://api.telebot.org"

type Bot interface {
	GetUpdates(UpdatesRequest) (*UpdateResponse, error)
	SendRequest(Request) (*http.Response, error)
	SendMessage(Request) (*MessageResponse, error)
}

type LongPoller interface {
	Run() error
}

func NewSimpleBot(Token string, client HttpClient) (SimpleBot, error) {
	return SimpleBot{
		apiEndpoint: DefaultApiUrl,
		client:      client,
		token:       Token,
		chatStates:  make(map[int]interface{}),
	}, nil
}

type SimpleBot struct {
	apiEndpoint string
	token       string
	client      HttpClient
	chatStates  map[int]interface{}
}

func (tb SimpleBot) GetUpdates(req UpdatesRequest) (resp *UpdateResponse, err error) {
	var httpResp *http.Response
	if httpResp, err = tb.SendRequest(req); err == nil {
		defer httpResp.Body.Close()
		resp = &UpdateResponse{}
		err = resp.Parse(httpResp.Body)
	}
	return
}

func (tb SimpleBot) SendRequest(botReq Request) (resp *http.Response, err error) {
	values, method, err := botReq.GetParams()
	if err == nil {
		var httpReq *http.Request
		url := fmt.Sprintf("%s/bot%s/%s", tb.apiEndpoint, tb.token, method)
		httpReq, err = http.NewRequest("POST", url, strings.NewReader(values.Encode()))
		if err == nil {
			httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp, err = tb.client.Do(httpReq)
		}
	}
	return
}

func (tb SimpleBot) SendMessage(req Request) (resp *MessageResponse, err error) {
	var httpResp *http.Response
	if httpResp, err = tb.SendRequest(req); err == nil {
		defer httpResp.Body.Close()
		resp = &MessageResponse{}
		err = resp.Parse(httpResp.Body)
	}
	return
}

type ContextPoller struct {
	bot            Bot
	offset         int
	logger         Logger
	mutex          sync.Mutex
	updateHandlers []UpdateHandler
}

func (lp *ContextPoller) getUpdates(ch chan UpdateResponse) (resp *UpdateResponse, err error) {
	lp.mutex.Lock()
	defer lp.mutex.Unlock()
	resp, err = lp.bot.GetUpdates(UpdatesRequest{Offset: lp.offset})
	if len(resp.Result) > 0 {
		lp.offset = resp.Result[len(resp.Result)-1].UpdateId + 1
	}

	if ch != nil {
		ch <- *resp
	}
	return
}

func (lp *ContextPoller) proceedUpdates(ctx context.Context) (err error) {
	chResp := make(chan UpdateResponse)
	go lp.getUpdates(chResp)

	select {
	case <-ctx.Done():
		lp.logger.infoF("SimplePoller proceed update closed by context")
		return nil

	case resp := <-chResp:
		if err != nil {
			lp.logger.errorF("SimplePoller proceed update error '%s'", err.Error())

			for _, update := range resp.Result {
				for _, handler := range lp.updateHandlers {
					ch := make(chan error)
					ctxWithCancel, cancelCtx := context.WithCancel(ctx)
					defer cancelCtx()

					go handler.ProceedUpdate(ctxWithCancel, lp.bot, update, ch)

					select {
					case <-ctx.Done():
						lp.logger.infoF("SimplePoller proceed update closed by context")
						return nil

					case err := <-ch:
						if err != nil {
							lp.logger.errorF("SimplePoller proceed update error '%s'", err.Error())
						}
					}
				}
			}
		}
	}
	return nil
}

type BaseLongPoller struct {
	*ContextPoller
}

func NewBaseLongPoller(tb Bot, handlers ...UpdateHandler) BaseLongPoller {
	ctxPoller := ContextPoller{bot: tb, updateHandlers: handlers}
	return BaseLongPoller{ContextPoller: &ctxPoller}
}

func (lp *BaseLongPoller) Run(ctx context.Context) error {
	for {
		if err := lp.proceedUpdates(ctx); err != nil {
			lp.logger.errorF("SimplePoller proceed update error '%s'", err.Error())
		}
	}
}
