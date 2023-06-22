package telebot

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type botMock struct {
	request       Request
	updateRequest UpdatesRequest
	err           error
}

func (bm *botMock) GetUpdates(ctx context.Context, ur UpdatesRequest) (UpdateResponse, error) {
	bm.updateRequest = ur
	return UpdateResponse{}, bm.err
}

func (bm *botMock) Send(ctx context.Context, r Request) (MessageResponse, error) {
	bm.request = r
	return MessageResponse{}, bm.err
}

func TestMessageGetCommand(t *testing.T) {
	tests := map[string]struct {
		text string
		want string
	}{
		"Text without command": {
			text: "some text",
			want: "",
		},
		"Command": {
			text: "/start",
			want: "start",
		},
		"Command with bot name": {
			text: "/start@bot",
			want: "start",
		},
		"Command with text": {
			text: "/start some text",
			want: "start",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			msg := Message{Text: test.text}
			if msg.GetCommand() != test.want {
				t.Fail()
			}
		})
	}
}

func TestMessageIsCommand(t *testing.T) {
	tests := map[string]struct {
		text string
		want bool
	}{
		"Text without command": {
			text: "some text",
			want: false,
		},
		"Command": {
			text: "/start",
			want: true,
		},
		"Command with bot name": {
			text: "/start@bot",
			want: true,
		},
		"Command with text": {
			text: "/start some text",
			want: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			msg := Message{Text: test.text}
			if msg.IsCommand() != test.want {
				t.Fail()
			}
		})
	}
}

func TestMessage_DeleteMessage(t *testing.T) {
	bm := botMock{}
	want := DeleteMessageRequest{ChatId: 1, MessageId: 10}
	_, err := Message{Chat: Chat{Id: 1}, MessageId: 10}.DeleteMessage(context.Background(), &bm)
	if err != nil {
		t.Errorf("Message.DeleteMessage() error = %v, wantErr %v", err, nil)
		return
	}
	if diff := cmp.Diff(bm.request, want); diff != "" {
		t.Errorf("Message.DeleteMessage() difference: %v", diff)
	}
}

func TestMessage_Edit(t *testing.T) {
	bm := botMock{}
	text := "New text"
	kbd := [][]InlineKeyboardButton{{{Text: "Button"}}}
	want := EditMessageTextRequest{ChatId: 1, MessageId: 10, Text: text, ReplyMarkup: kbd}
	_, err := Message{Chat: Chat{Id: 1}, MessageId: 10, Text: text, ReplyMarkup: kbd}.Edit(context.Background(), &bm)
	if err != nil {
		t.Errorf("Message.Edit() error = %v, wantErr %v", err, nil)
		return
	}
	if diff := cmp.Diff(bm.request, want); diff != "" {
		t.Errorf("Message.Edit() difference: %v", diff)
	}
}

func TestMessage_EditKeyboad(t *testing.T) {
	bm := botMock{}
	kbd := InlineKeyboardMarkup{[][]InlineKeyboardButton{{{Text: "Button"}}}}
	want := EditMessageReplyMarkup{ChatId: 1, MessageId: 10, ReplyMarkup: kbd}
	_, err := Message{Chat: Chat{Id: 1}, MessageId: 10}.EditKeyboard(context.Background(), &bm, kbd)
	if err != nil {
		t.Errorf("Message.EditKeyboard() error = %v, wantErr %v", err, nil)
		return
	}
	if diff := cmp.Diff(bm.request, want); diff != "" {
		t.Errorf("Message.EditKeyboard() difference: %v", diff)
	}
}

func TestMessage_EditMR(t *testing.T) {
	bm := botMock{}
	text := "New text"
	kbd := [][]InlineKeyboardButton{{{Text: "Button"}}}
	emr := EditMessageTextRequest{Text: text, ReplyMarkup: kbd}
	want := EditMessageTextRequest{ChatId: 1, MessageId: 10, Text: text, ReplyMarkup: kbd}
	_, err := Message{Chat: Chat{Id: 1}, MessageId: 10, Text: "Old text"}.EditMR(context.Background(), &bm, emr)
	if err != nil {
		t.Errorf("Message.EditText() error = %v, wantErr %v", err, nil)
		return
	}
	if diff := cmp.Diff(bm.request, want); diff != "" {
		t.Errorf("Message.EditText() difference: %v", diff)
	}
}

func TestMessage_EditText(t *testing.T) {
	bm := botMock{}
	text := "New text"
	kbd := [][]InlineKeyboardButton{{{Text: "Button"}}}
	want := EditMessageTextRequest{ChatId: 1, MessageId: 10, Text: text}
	_, err := Message{Chat: Chat{Id: 1}, MessageId: 10, ReplyMarkup: kbd}.EditText(context.Background(), &bm, text)
	if err != nil {
		t.Errorf("Message.EditText() error = %v, wantErr %v", err, nil)
		return
	}
	if diff := cmp.Diff(bm.request, want); diff != "" {
		t.Errorf("Message.EditText() difference: %v", diff)
	}
}

func TestMessage_ReplyText(t *testing.T) {
	bm := botMock{}
	text := "Reply text"
	want := MessageRequest{ChatId: 1, ReplyToMessageId: 10, Text: text}
	_, err := Message{Chat: Chat{Id: 1}, MessageId: 10}.ReplyText(context.Background(), &bm, text)
	if err != nil {
		t.Errorf("Message.ReplyText() error = %v, wantErr %v", err, nil)
		return
	}
	if diff := cmp.Diff(bm.request, want); diff != "" {
		t.Errorf("Message.ReplyText() difference: %v", diff)
	}
}

func TestMessage_ReplyMessage(t *testing.T) {
	bm := botMock{}
	text := "Send text"
	mr := MessageRequest{Text: text, DisableNotification: true}
	want := MessageRequest{ChatId: 1, Text: text, DisableNotification: true}
	_, err := Message{Chat: Chat{Id: 1}, Text: text}.ReplyMR(context.Background(), &bm, mr)
	if err != nil {
		t.Errorf("Message.Send() error = %v, wantErr %v", err, nil)
		return
	}
	if diff := cmp.Diff(bm.request, want); diff != "" {
		t.Errorf("Message.Send() difference: %v", diff)
	}
}

func TestMessage_Send(t *testing.T) {
	bm := botMock{}
	text := "Send text"
	kbd := [][]InlineKeyboardButton{{{Text: "Button"}}}
	want := MessageRequest{ChatId: 1, Text: text, ReplyMarkup: kbd}
	_, err := Message{Chat: Chat{Id: 1}, Text: text, ReplyMarkup: kbd}.Send(context.Background(), &bm)
	if err != nil {
		t.Errorf("Message.Send() error = %v, wantErr %v", err, nil)
		return
	}
	if diff := cmp.Diff(bm.request, want); diff != "" {
		t.Errorf("Message.Send() difference: %v", diff)
	}
}

func TestMessage_SendText(t *testing.T) {
	bm := botMock{}
	text := "Reply text"
	kbd := [][]InlineKeyboardButton{{{Text: "Button"}}}
	want := MessageRequest{ChatId: 1, Text: text} // Send only text message to same chat
	_, err := Message{Chat: Chat{Id: 1}, MessageId: 10, ReplyMarkup: kbd}.SendText(context.Background(), &bm, text)
	if err != nil {
		t.Errorf("Message.SendText() error = %v, wantErr %v", err, nil)
		return
	}
	if diff := cmp.Diff(bm.request, want); diff != "" {
		t.Errorf("Message.SendText() difference: %v", diff)
	}
}

func TestCallbackQuery_Answer(t *testing.T) {
	bm := botMock{}
	text := "Answer text"
	want := AnswerCallbackQueryRequest{CallbackQueryId: "1", Text: text}
	_, err := CallbackQuery{Id: "1"}.Answer(context.Background(), &bm, text)
	if err != nil {
		t.Errorf("CallbackQuery.Answer() error = %v, wantErr %v", err, nil)
		return
	}
	if diff := cmp.Diff(bm.request, want); diff != "" {
		t.Errorf("CallbackQuery.Answer() difference: %v", diff)
	}
}
