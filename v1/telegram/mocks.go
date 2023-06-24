package telegram

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type BodyMock struct {
	Reader io.Reader
}

func (bm BodyMock) Read(p []byte) (n int, err error) {
	n, err = bm.Reader.Read(p)
	return
}

func (bm BodyMock) Close() error {
	return nil
}

type httpClientMock struct {
	body string
	err  error
}

func (hcm httpClientMock) Do(httpRequest *http.Request) (*http.Response, error) {
	httpResponse := http.Response{}
	httpResponse.Request = httpRequest
	httpResponse.Body = BodyMock{strings.NewReader(hcm.body)}
	return &httpResponse, hcm.err
}

type UpdateHandlerMock struct {
	err error
}

func (h UpdateHandlerMock) Proceed(ctx context.Context, tb Bot, u ...Update) error {
	return h.err
}

type requestMock struct {
	values url.Values
	method string
	err    error
}

func (rm requestMock) GetParams() (url.Values, string, error) {
	return rm.values, rm.method, rm.err
}
