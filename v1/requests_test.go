package telebot

import (
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMessageRequest_GetParams(t *testing.T) {
	wantMethod := "sendMessage"
	tests := map[string]struct {
		request    *MessageRequest
		want       url.Values
		wantMethod string
		wantErr    bool
	}{
		"Required fields": {
			request: &MessageRequest{
				ChatId: 586350636,
				Text:   "Example of text",
			},
			want: map[string][]string{
				"chat_id": {"586350636"},
				"text":    {"Example of text"},
			},
			wantMethod: wantMethod,
		},
		"Empty ChatId": {request: &MessageRequest{Text: "Example of text"}, wantErr: true},
		"Empty Text":   {request: &MessageRequest{ChatId: 586350636}, wantErr: true},
		"Empty Fields": {request: &MessageRequest{}, wantErr: true},
		"Fully filled fields": {
			request: &MessageRequest{
				ChatId:    586350636,
				Text:      "Example of text",
				ParseMode: "MarkdownV2",
				Entities: []MessageEntity{
					{
						Type:   "url",
						Offset: 0,
						Length: 5,
						Url:    "https://google.com",
					},
					{
						Type:   "mention",
						Offset: 6,
						Length: 5,
						User: &User{
							Id:        987654321,
							IsBot:     false,
							FirstName: "Firstname",
						},
					},
				},
				DisableWebPagePreview:    true,
				DisableNotification:      true,
				ReplyToMessageId:         1234,
				AllowSendingWithoutReply: true,
				ReplyMarkup: InlineKeyboardMarkup{
					InlineKeyboard: [][]InlineKeyboardButton{{
						{Text: "Button text 1", CallbackData: "Data1"},
						{Text: "Button text 2", CallbackData: "Data2"},
					}},
				},
			},
			want: map[string][]string{
				"allow_sending_without_reply": {"true"},
				"chat_id":                     {"586350636"},
				"disable_notification":        {"true"},
				"disable_web_page_preview":    {"true"},
				"entities":                    {`[{"type":"url","offset":0,"length":5,"url":"https://google.com"},{"type":"mention","offset":6,"length":5,"user":{"id":987654321,"is_bot":false,"first_name":"Firstname"}}]`},
				"parse_mode":                  {"MarkdownV2"},
				"reply_to_message_id":         {"1234"},
				"text":                        {"Example of text"},
				"reply_markup":                {`{"inline_keyboard":[[{"text":"Button text 1","callback_data":"Data1"},{"text":"Button text 2","callback_data":"Data2"}]]}`},
			},
			wantMethod: wantMethod,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			values, method, err := test.request.GetParams()

			if (err != nil) != test.wantErr {
				t.Errorf("MessageRequest.GetParams() error = %v", err)
				return
			}
			if diff := cmp.Diff(values, test.want); diff != "" {
				t.Errorf("MessageRequest.GetParams() difference %v", diff)
			}
			if method != test.wantMethod {
				t.Errorf("EditMessageTextRequest.GetParams() gotMethod = %v, want %v", method, test.wantMethod)
			}
		})
	}
}

func TestEditMessageTextRequest_GetParams(t *testing.T) {
	wantMethod := "editMessageText"
	tests := map[string]struct {
		request    *EditMessageTextRequest
		want       url.Values
		wantMethod string
		wantErr    bool
	}{
		"Chat Message parameters": {
			request: &EditMessageTextRequest{
				ChatId:    10,
				MessageId: 100,
				Text:      "Example of text",
			},
			want: map[string][]string{
				"chat_id":    {"10"},
				"message_id": {"100"},
				"text":       {"Example of text"},
			},
			wantMethod: wantMethod,
		},
		"Inline Message parameters": {
			request: &EditMessageTextRequest{
				InlineMessageId: "20",
				Text:            "Example of text",
			},
			want: map[string][]string{
				"inline_message_id": {"20"},
				"text":              {"Example of text"},
			},
			wantMethod: wantMethod,
		},
		"Error fields": {
			request: &EditMessageTextRequest{
				Text: "Example of text",
			},
			wantErr: true,
		},
		"Fully filled parameters": {
			request: &EditMessageTextRequest{
				ChatId:    10,
				MessageId: 100,
				Text:      "Example of text",
				ParseMode: "MarkdownV2",
				Entities: []MessageEntity{
					{
						Type:   "url",
						Offset: 0,
						Length: 5,
						Url:    "https://google.com",
					},
					{
						Type:   "mention",
						Offset: 6,
						Length: 5,
						User: &User{
							Id:        987654321,
							IsBot:     false,
							FirstName: "Firstname",
						},
					},
				},
				DisableWebPagePreview: true,
				ReplyMarkup: InlineKeyboardMarkup{
					InlineKeyboard: [][]InlineKeyboardButton{{
						{Text: "Button text 1", CallbackData: "Data1"},
						{Text: "Button text 2", CallbackData: "Data2"},
					}},
				},
			},
			want: map[string][]string{
				"chat_id":                  {"10"},
				"message_id":               {"100"},
				"disable_web_page_preview": {"true"},
				"entities":                 {`[{"type":"url","offset":0,"length":5,"url":"https://google.com"},{"type":"mention","offset":6,"length":5,"user":{"id":987654321,"is_bot":false,"first_name":"Firstname"}}]`},
				"parse_mode":               {"MarkdownV2"},
				"text":                     {"Example of text"},
				"reply_markup":             {`{"inline_keyboard":[[{"text":"Button text 1","callback_data":"Data1"},{"text":"Button text 2","callback_data":"Data2"}]]}`},
			},
			wantMethod: wantMethod,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			values, method, err := test.request.GetParams()
			if (err != nil) != test.wantErr {
				t.Errorf("EditMessageTextRequest.GetParams() error = %v", err)
				return
			}
			if diff := cmp.Diff(values, test.want); diff != "" {
				t.Errorf("EditMessageTextRequest.GetParams() difference %v", diff)
			}
			if method != test.wantMethod {
				t.Errorf("EditMessageTextRequest.GetParams() gotMethod = %v, want %v", method, test.wantMethod)
			}
		})
	}
}

func TestSetMyCommandsRequest_GetParams(t *testing.T) {
	tests := map[string]struct {
		request *SetMyCommandsRequest
		want    map[string]string
	}{
		"Commands without Scope": {
			request: &SetMyCommandsRequest{
				Commands: []BotCommand{
					{Command: "start", Description: "Start description"},
					{Command: "help", Description: "Help description"},
				},
			},
			want: map[string]string{
				"commands": `[{"command":"start","description":"Start description"},{"command":"help","description":"Help description"}]`,
			},
		},
		"Commands with BotCommandScopeAllPrivateChats": {
			request: &SetMyCommandsRequest{
				Commands: []BotCommand{
					{Command: "start", Description: "Start description"},
					{Command: "help", Description: "Help description"},
				},
				Scope: BotCommandScope{Type: "all_private_chats"},
			},
			want: map[string]string{
				"commands": `[{"command":"start","description":"Start description"},{"command":"help","description":"Help description"}]`,
				"scope":    `{"type":"all_private_chats"}`,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			values, method, err := test.request.GetParams()

			if err != nil {
				t.Fail()
			}

			if method != "setMyCommands" {
				t.Fail()
			}

			for name, val := range test.want {
				valStr := values.Get(name)
				if valStr != val {
					t.Fail()
				}
			}
		})
	}
}

func TestDeleteMessageRequest_GetParams(t *testing.T) {
	tests := map[string]struct {
		request *DeleteMessageRequest
		want    map[string]string
	}{
		"Commands without Scope": {
			request: &DeleteMessageRequest{ChatId: 12345, MessageId: 54321},
			want:    map[string]string{"chat_id": "12345", "message_id": "54321"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			values, method, err := test.request.GetParams()

			if err != nil {
				t.Fail()
			}

			if method != "deleteMessage" {
				t.Fail()
			}

			for name, val := range test.want {
				valStr := values.Get(name)
				if valStr != val {
					t.Fail()
				}
			}
		})
	}
}

func TestInvoiceRequest_GetParams(t *testing.T) {
	wantMethod := "sendInvoice"
	tests := map[string]struct {
		request    *InvoiceRequest
		want       url.Values
		wantMethod string
		wantErr    bool
	}{
		"Required fields": {
			request: &InvoiceRequest{
				ChatId:        10,
				Title:         "Test invoice",
				Description:   "Test Description",
				Payload:       "Test pyload",
				ProviderToken: "PAY_TOKEN",
				Currency:      "RUB",
				Prices:        []LabeledPrice{{Label: "GOOD", Amount: 10}},
			},
			want: map[string][]string{
				"chat_id":        {"10"},
				"title":          {"Test invoice"},
				"description":    {"Test Description"},
				"payload":        {"Test pyload"},
				"provider_token": {"PAY_TOKEN"},
				"currency":       {"RUB"},
				"prices":         {`[{"label":"GOOD","amount":10}]`},
			},
			wantMethod: wantMethod,
		},
		"With keyboard": {
			request: &InvoiceRequest{
				ChatId:        10,
				Title:         "Test invoice",
				Description:   "Test Description",
				Payload:       "Test pyload",
				ProviderToken: "PAY_TOKEN",
				Currency:      "RUB",
				Prices:        []LabeledPrice{{Label: "GOOD", Amount: 10}},
				ReplyMarkup:   InlineKeyboardMarkup{[][]InlineKeyboardButton{{{Text: "Button"}}}},
			},
			want: map[string][]string{
				"chat_id":        {"10"},
				"title":          {"Test invoice"},
				"description":    {"Test Description"},
				"payload":        {"Test pyload"},
				"provider_token": {"PAY_TOKEN"},
				"currency":       {"RUB"},
				"prices":         {`[{"label":"GOOD","amount":10}]`},
				"reply_markup":   {`{"inline_keyboard":[[{"text":"Button"}]]}`},
			},
			wantMethod: wantMethod,
		},
		"Invalid fields": {
			request: &InvoiceRequest{
				ChatId:        10,
				Title:         "Test invoice",
				ProviderToken: "PAY_TOKEN",
				Currency:      "RUB",
				Prices:        []LabeledPrice{{Label: "GOOD", Amount: 10}},
			},
			wantErr: true,
		},
		"Empty fields": {request: &InvoiceRequest{}, wantErr: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			values, method, err := test.request.GetParams()
			if (err != nil) != test.wantErr {
				t.Errorf("InvoiceRequest.GetParams() error = %v", err)
				return
			}
			if diff := cmp.Diff(values, test.want); diff != "" {
				t.Errorf("InvoiceRequest.GetParams() difference %v", diff)
			}
			if method != test.wantMethod {
				t.Errorf("InvoiceRequest.GetParams() gotMethod = %v, want %v", method, test.wantMethod)
			}
		})
	}
}

func TestEditMessageReplyMarkup_GetParams(t *testing.T) {
	wantMethod := "editMessageReplyMarkup"
	type fields struct {
		ChatId          interface{}
		MessageId       int
		InlineMessageId string
		ReplyMarkup     InlineKeyboardMarkup
	}
	tests := []struct {
		name       string
		fields     fields
		wantVal    url.Values
		wantMethod string
		wantErr    bool
	}{
		{
			name: "Required fields",
			fields: fields{
				ChatId:    10,
				MessageId: 100,
			},
			wantVal:    map[string][]string{"chat_id": {"10"}, "message_id": {"100"}},
			wantMethod: wantMethod,
		},
		{
			name: "With keyboard",
			fields: fields{
				ChatId:      10,
				MessageId:   100,
				ReplyMarkup: InlineKeyboardMarkup{[][]InlineKeyboardButton{{{Text: "Button"}}}},
			},
			wantVal:    map[string][]string{"chat_id": {"10"}, "message_id": {"100"}, "reply_markup": {`{"inline_keyboard":[[{"text":"Button"}]]}`}},
			wantMethod: wantMethod,
		},
		{
			name: "Inline message",
			fields: fields{
				InlineMessageId: "20",
			},
			wantVal:    map[string][]string{"inline_message_id": {"20"}},
			wantMethod: wantMethod,
		},
		{
			name:    "Required fields error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := EditMessageReplyMarkup{
				ChatId:          tt.fields.ChatId,
				MessageId:       tt.fields.MessageId,
				InlineMessageId: tt.fields.InlineMessageId,
				ReplyMarkup:     tt.fields.ReplyMarkup,
			}
			gotVal, gotMethod, err := req.GetParams()
			if (err != nil) != tt.wantErr {
				t.Errorf("EditMessageReplyMarkup.GetParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(gotVal, tt.wantVal); diff != "" {
				t.Errorf("EditMessageReplyMarkup.GetParams() difference %v", diff)
			}
			if gotMethod != tt.wantMethod {
				t.Errorf("EditMessageReplyMarkup.GetParams() gotMethod = %v, want %v", gotMethod, tt.wantMethod)
			}
		})
	}
}

func TestAnswerCallbackQueryRequest_GetParams(t *testing.T) {
	wantMethod := "answerCallbackQuery"
	type fields struct {
		CallbackQueryId string
		Text            string
		ShowAlert       bool
		URL             string
		CacheTime       int
	}
	tests := []struct {
		name       string
		fields     fields
		wantVal    url.Values
		wantMethod string
		wantErr    bool
	}{
		{
			name:       "Required fields",
			fields:     fields{CallbackQueryId: "1010"},
			wantVal:    map[string][]string{"callback_query_id": {"1010"}, "show_alert": {"false"}},
			wantMethod: wantMethod,
		},
		{
			name:    "Empty fields",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "Full fields",
			fields: fields{
				CallbackQueryId: "1010",
				Text:            "Ok",
				ShowAlert:       true,
				URL:             "http:/url.local",
				CacheTime:       100,
			},
			wantVal: map[string][]string{
				"callback_query_id": {"1010"},
				"text":              {"Ok"},
				"show_alert":        {"true"},
				"url":               {"http:/url.local"},
				"cache_time":        {"100"},
			},
			wantMethod: wantMethod,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := AnswerCallbackQueryRequest{
				CallbackQueryId: tt.fields.CallbackQueryId,
				Text:            tt.fields.Text,
				ShowAlert:       tt.fields.ShowAlert,
				URL:             tt.fields.URL,
				CacheTime:       tt.fields.CacheTime,
			}
			gotVal, gotMethod, err := req.GetParams()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnswerCallbackQueryRequest.GetParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(gotVal, tt.wantVal); diff != "" {
				t.Errorf("AnswerCallbackQueryRequest.GetParams() difference %v", diff)
			}
			if gotMethod != tt.wantMethod {
				t.Errorf("AnswerCallbackQueryRequest.GetParams() gotMethod = %v, want %v", gotMethod, tt.wantMethod)
			}
		})
	}
}
