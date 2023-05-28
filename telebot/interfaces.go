package telebot

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type UpdateHandler interface {
	ProceedUpdate(context.Context, Bot, Update, chan error) error
}

type CallbackHandler interface {
	ProceedCallback(*CallbackQuery) error
}

type MessageHandler interface {
	ProceedMessage(tm *Message) error
}

type MessageRequestHelper interface {
	GetEditMR() EditMessageTextRequest
	GetMR() MessageRequest
}

type CallbackDataParser interface {
	GetAction() string
	GetPrefix() string
	GetState() State
	GetValue() string
	Parse(string) error
	SetState(state State)
}

type KeyboardHelper interface {
	GetKeyboard() interface{}
	GetText() string
}

type Request interface {
	GetParams() (url.Values, string, error)
}

type Response interface {
	Parse(reader io.Reader) error
}

type StateBuilder interface {
	GetStateProvider(State) (StateProvider, error)
}

type StateProvider interface {
	GetRequests() []StateRequest
	Proceed() (State, error)
}

type StateRepository interface {
	Get(ChatId int) ([]State, error)
	GetByData(Data string) ([]State, error)
	GetByMessage(msg Message) (State, error)
	Set(State) error
	Clear(State) error
}
