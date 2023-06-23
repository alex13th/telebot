package telebot

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"sync"
)

type MessageRequest struct {
	ChatId                   interface{}     `json:"chat_id"`
	Text                     string          `json:"text"`
	ParseMode                string          `json:"parse_mode"`
	Entities                 []MessageEntity `json:"entities"`
	DisableWebPagePreview    bool            `json:"disable_web_page_preview"`
	DisableNotification      bool            `json:"disable_notification"`
	ReplyToMessageId         int             `json:"reply_to_message_id"`
	AllowSendingWithoutReply bool            `json:"allow_sending_without_reply"`
	ReplyMarkup              interface{}     `json:"reply_markup"`
}

func (req MessageRequest) GetParams() (val url.Values, method string, err error) {
	method = "sendMessage"
	if req.ChatId == nil || req.Text == "" {
		return nil, "",
			fmt.Errorf("required fields not defined, ChatId: %v, Text: %s", req.ChatId, req.Text)
	}

	val = url.Values{}
	val.Add("chat_id", fmt.Sprint(req.ChatId))
	val.Add("text", req.Text)
	if req.ParseMode != "" {
		val.Add("parse_mode", req.ParseMode)
	}
	if req.DisableWebPagePreview {
		val.Add("disable_web_page_preview", strconv.FormatBool(req.DisableWebPagePreview))
	}
	if req.DisableNotification {
		val.Add("disable_notification", strconv.FormatBool(req.DisableNotification))
	}
	if req.ReplyToMessageId > 0 {
		val.Add("reply_to_message_id", strconv.Itoa(req.ReplyToMessageId))
	}
	if req.AllowSendingWithoutReply {
		val.Add("allow_sending_without_reply", strconv.FormatBool(req.AllowSendingWithoutReply))
	}
	if len(req.Entities) > 0 {
		data, err := json.Marshal(req.Entities)
		if err != nil {
			return nil, "", err
		}
		val.Add("entities", string(data))
	}

	if req.ReplyMarkup != nil {
		data, err := json.Marshal(req.ReplyMarkup)
		if err != nil {
			return nil, "", err
		}
		val.Add("reply_markup", string(data))
	}
	return
}

type EditMessageTextRequest struct {
	ChatId                interface{}     `json:"chat_id,omitempty"`
	MessageId             int             `json:"message_id,omitempty"`
	InlineMessageId       string          `json:"inline_message_id,omitempty"`
	Text                  string          `json:"text"`
	ParseMode             string          `json:"parse_mode"`
	Entities              []MessageEntity `json:"entities"`
	DisableWebPagePreview bool            `json:"disable_web_page_preview"`
	ReplyMarkup           interface{}     `json:"reply_markup"`
}

func (req EditMessageTextRequest) GetParams() (val url.Values, method string, err error) {
	method = "editMessageText"
	if (req.ChatId == nil || req.MessageId == 0) && (req.InlineMessageId == "") {
		return nil, "",
			fmt.Errorf("required fields not defined, ChatId: %v, MessageId: %d, InlineMessageId: %s", req.ChatId, req.MessageId, req.InlineMessageId)
	}

	val = url.Values{}

	if req.ChatId != nil {
		val.Add("chat_id", fmt.Sprint(req.ChatId))
	}

	if req.MessageId != 0 {
		val.Add("message_id", strconv.Itoa(req.MessageId))
	}

	if req.InlineMessageId != "" {
		val.Add("inline_message_id", req.InlineMessageId)
	}

	val.Add("text", req.Text)
	if req.ParseMode != "" {
		val.Add("parse_mode", req.ParseMode)
	}
	if req.DisableWebPagePreview {
		val.Add("disable_web_page_preview", strconv.FormatBool(req.DisableWebPagePreview))
	}
	if len(req.Entities) > 0 {
		data, err := json.Marshal(req.Entities)
		if err != nil {
			return nil, "", err
		}
		val.Add("entities", string(data))
	}

	if req.ReplyMarkup != nil {
		data, err := json.Marshal(req.ReplyMarkup)
		if err != nil {
			return nil, "", err
		}
		val.Add("reply_markup", string(data))
	}
	return
}

type EditMessageReplyMarkup struct {
	ChatId          interface{}          `json:"chat_id"`
	MessageId       int                  `json:"message_id"`
	InlineMessageId string               `json:"inline_message_id"`
	ReplyMarkup     InlineKeyboardMarkup `json:"reply_markup"`
}

func (req EditMessageReplyMarkup) GetParams() (val url.Values, method string, err error) {
	method = "editMessageReplyMarkup"

	if (req.ChatId == nil || req.MessageId == 0) && (req.InlineMessageId == "") {
		return nil, "",
			fmt.Errorf("required fields not defined, ChatId: %v, MessageId: %d, InlineMessageId: %s", req.ChatId, req.MessageId, req.InlineMessageId)
	}

	val = url.Values{}

	if req.ChatId != nil {
		val.Add("chat_id", fmt.Sprint(req.ChatId))
	}

	if req.MessageId != 0 {
		val.Add("message_id", strconv.Itoa(req.MessageId))
	}

	if req.InlineMessageId != "" {
		val.Add("inline_message_id", req.InlineMessageId)
	}

	if req.ReplyMarkup.InlineKeyboard != nil {
		data, err := json.Marshal(req.ReplyMarkup)
		if err != nil {
			return nil, "", err
		}
		val.Add("reply_markup", string(data))
	}
	return
}

type AnswerCallbackQueryRequest struct {
	CallbackQueryId string `json:"callback_query_id"`
	Text            string `json:"text"`
	ShowAlert       bool   `json:"show_alert"`
	URL             string `json:"url"`
	CacheTime       int    `json:"cache_time"`
}

func (req AnswerCallbackQueryRequest) GetParams() (val url.Values, method string, err error) {
	method = "answerCallbackQuery"
	if req.CallbackQueryId == "" {
		return nil, "",
			fmt.Errorf("required fields not defined, CallbackQueryId: %s", req.CallbackQueryId)
	}
	val = url.Values{}
	val.Add("callback_query_id", req.CallbackQueryId)

	if req.Text != "" {
		val.Add("text", req.Text)
	}
	val.Add("show_alert", strconv.FormatBool(req.ShowAlert))
	if req.URL != "" {
		val.Add("url", req.URL)
	}
	if req.CacheTime != 0 {
		val.Add("cache_time", strconv.Itoa(req.CacheTime))
	}
	return
}

type SetMyCommandsRequest struct {
	mu           sync.RWMutex
	Commands     []BotCommand `json:"commands"`
	Scope        interface{}  `json:"scope"`
	LanguageCode string       `json:"language_code"`
}

func (req *SetMyCommandsRequest) GetParams() (val url.Values, method string, err error) {
	method = "setMyCommands"
	val = url.Values{}
	req.mu.RLock()
	data, err := json.Marshal(req.Commands)
	if err != nil {
		return nil, "", err
	}
	val.Add("commands", string(data))

	if req.Scope != nil {
		data, err := json.Marshal(req.Scope)
		if err != nil {
			return nil, "", err
		}
		val.Add("scope", string(data))
	}

	val.Add("language_code", req.LanguageCode)
	defer req.mu.RUnlock()
	return
}

type DeleteMessageRequest struct {
	ChatId    interface{} `json:"chat_id"`
	MessageId int         `json:"message_id"`
}

func (req DeleteMessageRequest) GetParams() (url.Values, string, error) {
	method := "deleteMessage"
	val := url.Values{}
	val.Add("chat_id", fmt.Sprint(req.ChatId))
	val.Add("message_id", fmt.Sprint(req.MessageId))
	return val, method, nil
}

type LabeledPrice struct {
	Label  interface{} `json:"label"`
	Amount interface{} `json:"amount"`
}
type InvoiceRequest struct {
	ChatId        interface{}          `json:"chat_id"`
	Title         string               `json:"title"`
	Description   string               `json:"description"`
	Payload       string               `json:"payload"`
	ProviderToken string               `json:"provider_token"`
	Currency      string               `json:"currency"`
	Prices        []LabeledPrice       `json:"prices"`
	ReplyMarkup   InlineKeyboardMarkup `json:"reply_markup"`
}

func (req InvoiceRequest) GetParams() (val url.Values, method string, err error) {
	method = "sendInvoice"
	if req.ChatId == nil || req.Title == "" || req.Description == "" ||
		req.Payload == "" || req.ProviderToken == "" || req.Currency == "" ||
		req.Prices == nil {

		return nil, "",
			fmt.Errorf("required fields not defined, ChatId: %v, Title: %s, Description: %s, "+
				"Payload: %s, ProviderToken: %s, Prices: %v",
				req.ChatId, req.Title, req.Description, req.Payload, req.ProviderToken, req.Prices)
	}

	val = url.Values{}
	val.Add("chat_id", fmt.Sprint(req.ChatId))
	val.Add("title", req.Title)
	val.Add("description", req.Description)
	val.Add("payload", req.Payload)
	val.Add("provider_token", req.ProviderToken)
	val.Add("currency", req.Currency)

	data, err := json.Marshal(req.Prices)
	if err != nil {
		return nil, "", err
	}
	val.Add("prices", string(data))

	if req.ReplyMarkup.InlineKeyboard != nil {
		data, err := json.Marshal(req.ReplyMarkup)
		if err != nil {
			return nil, "", err
		}
		val.Add("reply_markup", string(data))
	}
	return
}
