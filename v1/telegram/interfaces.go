package telegram

import (
	"context"
)

type Bot interface {
	GetUpdates(context.Context, UpdatesRequest) (UpdateResponse, error)
	Send(context.Context, Request) (MessageResponse, error)
}

type UpdateHandler interface {
	Proceed(context.Context, Bot, ...Update) error
}

type CallbackHandler interface {
	ProceedCallback(CallbackQuery) error
}

type MessageHandler interface {
	ProceedMessage(tm Message) error
}
