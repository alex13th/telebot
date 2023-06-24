package telegram

import (
	"encoding/json"
	"fmt"
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
		var err error
		if ur.Result[i].Message, err = normalizeMessage(update.Message); err != nil {
			*ur = UpdateResponse{}
			return err
		}
		if ur.Result[i].EditedMessage, err = normalizeMessage(update.EditedMessage); err != nil {
			*ur = UpdateResponse{}
			return err
		}
		if ur.Result[i].ChannelPost, err = normalizeMessage(update.ChannelPost); err != nil {
			*ur = UpdateResponse{}
			return err
		}
		if ur.Result[i].EditedChannelPost, err = normalizeMessage(update.EditedChannelPost); err != nil {
			*ur = UpdateResponse{}
			return err
		}
	}
	return nil
}

func normalizeMessage(m Message) (Message, error) {
	if m.Chat.Id == nil {
		return m, nil
	}

	switch val := m.Chat.Id.(type) {
	case float64:
		m.Chat.Id = int(val)
	case int:
	case string:
	default:
		return Message{}, fmt.Errorf("invalid chat Id type %v", val)
	}

	if m.ReplyToMessage != nil {
		switch val := m.ReplyToMessage.Chat.Id.(type) {
		case float64:
			m.ReplyToMessage.Chat.Id = int(val)
		case int:
		case string:
		default:
			return Message{}, fmt.Errorf("invalid ReplyToMessage chat Id type %v", val)
		}
	}
	return m, nil
}

func (mr *MessageResponse) Parse(reader io.Reader) error {
	err := ParseJson(mr, reader)
	if err != nil {
		return err
	}
	mr.Result, err = normalizeMessage(mr.Result)

	if err != nil {
		*mr = MessageResponse{}
		return err
	}
	return nil
}
