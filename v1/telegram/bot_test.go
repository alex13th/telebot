package telegram

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSimpleBot_sendRequest(t *testing.T) {
	tb := NewSimpleBot("***Token***", httpClientMock{})
	tb.client = httpClientMock{}

	ctxErr := errors.New("context error")
	reqErr := errors.New("request error")

	data := []struct {
		name string
		req  requestMock
		ctx  context.Context
		err  error
	}{
		{
			name: "Valid",
			req: requestMock{
				method: "sendMessage",
				values: url.Values{"chat_id": {"586350636"}, "text": {"Message text"}},
			},
			ctx: context.Background(),
		},
		{
			name: "Without context",
			req: requestMock{
				method: "sendMessage",
				values: url.Values{"chat_id": {"586350636"}, "text": {"Message text"}},
			},
			err: ctxErr,
		},
		{
			name: "With error",
			req:  requestMock{err: reqErr},
			err:  reqErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			resp, err := tb.sendRequest(d.ctx, d.req)

			if err != nil {
				if !errors.Is(err, d.err) && d.err != ctxErr {
					t.Errorf("expected error '%v', but '%v'", d.err, err)
				}
				return
			}

			if resp.Request.Method != "POST" {
				t.Errorf("expected HTTP method %s, but %s", d.req.method, "POST")
			}

			url := fmt.Sprintf("%s/bot***Token***/%s", DefaultApiUrl, d.req.method)
			if resp.Request.URL.String() != url {
				t.Errorf("expected request URL %s, but %s", url, resp.Request.URL.String())
			}

			ct := "application/x-www-form-urlencoded"
			if len(resp.Request.Header["Content-Type"]) != 1 {
				t.Errorf("expected Content-Type count 1, but %d", len(resp.Request.Header["Content-Type"]))
			}
			if resp.Request.Header["Content-Type"][0] != ct {
				t.Errorf("expected Content-Type %s, but %s", ct, resp.Request.Header["Content-Type"][0])
			}

			bytes, err := io.ReadAll(resp.Request.Body)
			body := "chat_id=586350636&text=Message+text"
			if err != nil {
				t.Errorf("read request Body error %s", err)
			}
			if string(bytes) != body {
				t.Errorf("expected request Body %s, but %s", body, string(bytes))
			}
		})
	}
}

func TestSimpleBotSend(t *testing.T) {
	tb := NewSimpleBot(
		"***Token***",
		httpClientMock{
			body: `{
			"ok": true,
			"result": {
				"message_id": 2468,
				"from": {"id": 586350636,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
				"chat": {"id": 586350636,"first_name": "Alexey","last_name": "Sukharev","type": "private"},
				"date": 1630134810,
				"text": "Hello world!!!"
			}}`,
		})

	data := []struct {
		name string
		err  error
	}{
		{name: "Without error", err: nil},
		{name: "With error", err: errors.New("Mock error")},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			resp, err := tb.Send(context.Background(), requestMock{err: d.err})

			if err != nil {
				if !errors.Is(err, d.err) {
					t.Errorf("expected error '%s', but '%s'", d.err, err)
				}
				return
			}

			if !resp.Ok {
				t.Error("bot response not Ok")
			}

			mid := 2468
			if resp.Result.MessageId != mid {
				t.Errorf("expected message Id %d, but %d", mid, resp.Result.MessageId)
			}
		})
	}
}

func TestSimpleBotGetUpdates(t *testing.T) {
	httpErr := errors.New("HTTP error")
	data := []struct {
		name   string
		client httpClient
		err    error
	}{
		{
			name:   "Valid",
			client: httpClientMock{body: `{"ok": true,"result":[{"update_id":123130161},{"update_id":123130162},{"update_id":123130163}]}`},
		},
		{
			name:   "JSON error",
			client: httpClientMock{body: ""},
			err:    io.EOF,
		},
		{
			name:   "With HTTP error",
			client: httpClientMock{err: httpErr},
			err:    httpErr,
		},
		{
			name:   "With Telegram error",
			client: httpClientMock{body: `{"ok": false,"error_code":400,"description":"telegram API error"}`},
			err:    ErrStatus{ErrorCode: 400, Description: "telegram API error"},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			tb := NewSimpleBot("***Token***", d.client)

			resp, err := tb.GetUpdates(context.Background(), UpdatesRequest{})

			if err != nil {
				if !errors.Is(err, d.err) {
					t.Errorf("expected error '%v', but '%v'", d.err, err)
				}
				return
			}

			if !resp.Ok {
				t.Error("bot response not Ok")
			}

			updates := []Update{{UpdateId: 123130161}, {UpdateId: 123130162}, {UpdateId: 123130163}}
			if diff := cmp.Diff(resp.Result, updates); diff != "" {
				t.Errorf("expected update difference %s", diff)
			}
		})
	}
}

func TestSimplePoller_getUpdates(t *testing.T) {
	data := []struct {
		name string
		body string
		err  error
	}{
		{
			name: "Valid",
			body: `{"ok":true,"result":[{"update_id": 123130161},{"update_id": 123130162},{"update_id": 123130163}]}`,
		},
		{
			name: "HTTP client error",
			err:  errors.New("http client error"),
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			tb := NewSimpleBot(
				"***Token***",
				httpClientMock{body: d.body, err: d.err},
			)

			slp := NewSimplePoller(tb, nil)
			_, err := slp.getUpdates(context.Background())

			if !errors.Is(err, d.err) {
				t.Errorf("expected error '%s', but '%s'", d.err, err)
			}
		})
	}
}

func TestSimplePoller_proceedUpdates(t *testing.T) {
	body := `{"ok":true,"result":[{"update_id": 123130161},{"update_id": 123130162},{"update_id": 123130163}]}`
	data := []struct {
		name   string
		body   string
		offset int
		err    error
	}{
		{name: "Valid", body: body, offset: 123130164},
		{name: "Handler error", body: body, err: errors.New("update handler error"), offset: 0},
		{name: "Update error", body: "", offset: 0, err: io.EOF},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			tb := NewSimpleBot(
				"***Token***",
				httpClientMock{body: d.body},
			)

			pol := NewSimplePoller(tb, UpdateHandlerMock{err: d.err})
			offset, err := pol.ProceedUpdates(context.Background())

			if err != nil {
				if wrappedErr := errors.Unwrap(d.err); errors.Is(wrappedErr, d.err) {
					t.Errorf("expected error '%s', but '%v'", d.err, err)
				}
				return
			}

			if offset != d.offset {
				t.Errorf("expected offset is %d, but %d", 123130164, pol.offset)
			}
		})
	}
}

func Test_ErrStatus_Error(t *testing.T) {
	tests := []struct {
		name string
		se   ErrStatus
		want string
	}{
		{
			name: "Test errormessage",
			se:   ErrStatus{ErrorCode: 401, Description: "unauthorized access"},
			want: "telegram status error 'unauthorized access', error_code: 401",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.se.Error(); got != tt.want {
				t.Errorf(" ErrStatus.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdatesRequest_GetParams(t *testing.T) {
	wantMethod := "getUpdates"

	tests := []struct {
		name    string
		req     UpdatesRequest
		wantVal url.Values
		wantErr bool
	}{
		{name: "Simple", wantVal: url.Values{}},
		{name: "With offset", req: UpdatesRequest{Offset: 10}, wantVal: map[string][]string{"offset": {"10"}}},
		{
			name: "Fukk params",
			req:  UpdatesRequest{Offset: 10, Limit: 100, Timeout: 200, AllowedUpdates: []string{"message"}},
			wantVal: map[string][]string{
				"offset":          {"10"},
				"limit":           {"100"},
				"timeout":         {"200"},
				"allowed_updates": {"message"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotMethod, err := tt.req.GetParams()
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatesRequest.GetParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(gotVal, tt.wantVal); diff != "" {
				t.Errorf("UpdatesRequest.GetParams() gotVal difference: = %v", diff)
			}
			if gotMethod != wantMethod {
				t.Errorf("UpdatesRequest.GetParams() gotMethod = %v, want %v", gotMethod, wantMethod)
			}
		})
	}
}
