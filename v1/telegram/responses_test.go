package telegram

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_normalizeMessage(t *testing.T) {
	tests := []struct {
		name    string
		m       Message
		want    Message
		wantErr bool
	}{
		{name: "Chat is int", m: Message{Chat: Chat{Id: 10}}, want: Message{Chat: Chat{Id: 10}}},
		{name: "Chat is string", m: Message{Chat: Chat{Id: "@username"}}, want: Message{Chat: Chat{Id: "@username"}}},
		{name: "Chat is int float64", m: Message{Chat: Chat{Id: float64(123456789123456)}}, want: Message{Chat: Chat{Id: 123456789123456}}},
		{name: "Chat is fract float64", m: Message{Chat: Chat{Id: float64(1234567891.23456)}}, want: Message{Chat: Chat{Id: 1234567891}}},
		{name: "Chat is bool", m: Message{Chat: Chat{Id: true}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeMessage(tt.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("normalizeMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("normalizeMessage() difference: %s", diff)
			}
		})
		t.Run("ReplyToMessage - "+tt.name, func(t *testing.T) {
			m := Message{Chat: Chat{Id: 10}, ReplyToMessage: &tt.m}
			want := Message{}
			if !tt.wantErr {
				want = Message{Chat: Chat{Id: 10}, ReplyToMessage: &tt.want}
			}
			got, err := normalizeMessage(m)
			if (err != nil) != tt.wantErr {
				t.Errorf("normalizeMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("normalizeMessage() difference: %s", diff)
			}
		})
	}
}

func TestMessageResponse_Parse(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    MessageResponse
		wantErr bool
	}{
		{
			name: "Chat id is int",
			json: `{
				"ok": true,
				"result": {
					"message_id": 2468,
					"from": {"id": 10,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
					"chat": {"id": 1,"first_name": "Alexey","last_name": "Sukharev","type": "private"},
					"date": 1630134810,
					"text": "Hello world!!!"
				}
			}`,
			want: MessageResponse{
				Ok: true,
				Result: Message{
					MessageId: 2468,
					From:      User{Id: 10, IsBot: false, FirstName: "Alexey", LastName: "Sukharev", LanguageCode: "en"},
					Chat:      Chat{Id: 1, FirstName: "Alexey", LastName: "Sukharev", Type: "private"},
					Date:      1630134810,
					Text:      "Hello world!!!",
				},
			},
		},
		{
			name: "Chat id is string",
			json: `{
				"ok": true,
				"result": {
					"message_id": 2468,
					"from": {"id": 10,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
					"chat": {"id": "@username","first_name": "Alexey","last_name": "Sukharev","type": "private"},
					"date": 1630134810,
					"text": "Hello world!!!"
				}
			}`,
			want: MessageResponse{
				Ok: true,
				Result: Message{
					MessageId: 2468,
					From:      User{Id: 10, IsBot: false, FirstName: "Alexey", LastName: "Sukharev", LanguageCode: "en"},
					Chat:      Chat{Id: "@username", FirstName: "Alexey", LastName: "Sukharev", Type: "private"},
					Date:      1630134810,
					Text:      "Hello world!!!",
				},
			},
		},
		{
			name: "Wrong chat id",
			json: `{
				"ok": true,
				"result": {
					"message_id": 2468,
					"from": {"id": 10,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
					"chat": {"id": true,"first_name": "Alexey","last_name": "Sukharev","type": "private"},
					"date": 1630134810,
					"text": "Hello world!!!"
				}
			}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := &MessageResponse{}
			err := mr.Parse(strings.NewReader(tt.json))

			if (err != nil) != tt.wantErr {
				t.Errorf("MessageResponse.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(*mr, tt.want); diff != "" {
				t.Errorf("MessageResponse.Parse() difference: %s", diff)
			}
		})
	}
}

func TestUpdateResponse_Parse(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    UpdateResponse
		wantErr bool
	}{
		{
			name: "Two updates",
			json: `{
				"ok": true,
				"result": [
					{
						"update_id": 123130160,
						"message": {
							"message_id": 2468,
							"from": {"id": 10,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
							"chat": {"id": 1,"first_name": "Alexey","last_name": "Sukharev","type": "private"},
							"date": 1630134810,
							"text": "Hello world!!!"
						}
					},
					{
						"update_id": 123130161,
						"message": {
							"message_id": 2469,
							"from": {"id": 11,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
							"chat": {"id": 1,"first_name": "Alexey","last_name": "Sukharev","type": "private"},
							"date": 1630134810,
							"text": "Hello world!!!"
						}
					}
				]
			}`,
			want: UpdateResponse{
				Ok: true,
				Result: []Update{
					{
						UpdateId: 123130160,
						Message: Message{
							MessageId: 2468,
							From:      User{Id: 10, IsBot: false, FirstName: "Alexey", LastName: "Sukharev", LanguageCode: "en"},
							Chat:      Chat{Id: 1, FirstName: "Alexey", LastName: "Sukharev", Type: "private"},
							Date:      1630134810,
							Text:      "Hello world!!!",
						},
					},
					{
						UpdateId: 123130161,
						Message: Message{
							MessageId: 2469,
							From:      User{Id: 11, IsBot: false, FirstName: "Alexey", LastName: "Sukharev", LanguageCode: "en"},
							Chat:      Chat{Id: 1, FirstName: "Alexey", LastName: "Sukharev", Type: "private"},
							Date:      1630134810,
							Text:      "Hello world!!!",
						},
					},
				},
			},
		},
		{
			name: "Wrong message chat id",
			json: `{
				"ok": true,
				"result": [
					{
						"update_id": 123130160,
						"message": {
							"message_id": 2468,
							"from": {"id": 10,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
							"chat": {"id": true,"first_name": "Alexey","last_name": "Sukharev","type": "private"},
							"date": 1630134810,
							"text": "Hello world!!!"
						}
					}
				]
			}`,
			wantErr: true,
		},
		{
			name: "Wrong EditedMessage chat id",
			json: `{
				"ok": true,
				"result": [
					{
						"update_id": 123130160,
						"edited_message": {
							"message_id": 2468,
							"from": {"id": 10,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
							"chat": {"id": true,"first_name": "Alexey","last_name": "Sukharev","type": "private"},
							"date": 1630134810,
							"text": "Hello world!!!"
						}
					}
				]
			}`,
			wantErr: true,
		},
		{
			name: "Wrong ChannelPost chat id",
			json: `{
				"ok": true,
				"result": [
					{
						"update_id": 123130160,
						"channel_post": {
							"message_id": 2468,
							"from": {"id": 10,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
							"chat": {"id": true,"first_name": "Alexey","last_name": "Sukharev","type": "private"},
							"date": 1630134810,
							"text": "Hello world!!!"
						}
					}
				]
			}`,
			wantErr: true,
		},
		{
			name: "Wrong EditedChannelPost chat id",
			json: `{
				"ok": true,
				"result": [
					{
						"update_id": 123130160,
						"edited_channel_post": {
							"message_id": 2468,
							"from": {"id": 10,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
							"chat": {"id": true,"first_name": "Alexey","last_name": "Sukharev","type": "private"},
							"date": 1630134810,
							"text": "Hello world!!!"
						}
					}
				]
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := &UpdateResponse{}
			err := ur.Parse(strings.NewReader(tt.json))

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateResponse.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(*ur, tt.want); diff != "" {
				t.Errorf("UpdateResponse.Parse() difference: %s", diff)
			}
		})
	}
}
