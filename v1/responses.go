package telebot

import (
	"encoding/json"
	"io"
)

type UpdateResponse struct {
	Ok          bool               `json:"ok"`
	Result      []Update           `json:"result"`
	Description string             `json:"description"`
	ErrorCode   int                `json:"error_code"`
	Parameters  ResponseParameters `json:"parameters"`
}

type MessageResponse struct {
	Ok         bool               `json:"ok"`
	Result     Message            `json:"result"`
	ErrorCode  int                `json:"error_code"`
	Parameters ResponseParameters `json:"parameters"`
}

type ResponseParameters struct {
	MigrateToChatId int `json:"migrate_to_chat_id"`
	RetryAfter      int `json:"retry_after"`
}

func ParseJson(i interface{}, reader io.Reader) error {
	dec := json.NewDecoder(reader)
	return dec.Decode(i)
}

func (ur *UpdateResponse) Parse(reader io.Reader) error {
	if err := ParseJson(ur, reader); err != nil {
		return err
	}

	for i, update := range ur.Result {
		ur.Result[i].Message = normalizeMessage(update.Message)
		ur.Result[i].EditedMessage = normalizeMessage(update.EditedMessage)
		ur.Result[i].ChannelPost = normalizeMessage(update.ChannelPost)
		ur.Result[i].EditedChannelPost = normalizeMessage(update.EditedChannelPost)
	}
	return nil
}

func normalizeMessage(m Message) Message {
	switch val := m.Chat.Id.(type) {
	case float64:
		m.Chat.Id = int(val)
	}
	if m.ReplyToMessage != nil {
		switch val := m.ReplyToMessage.Chat.Id.(type) {
		case float64:
			m.ReplyToMessage.Chat.Id = int(val)
		}
	}
	return m
}

func (message *MessageResponse) Parse(reader io.Reader) error {
	return ParseJson(message, reader)
}
