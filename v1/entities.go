package telebot

import (
	"context"
	"regexp"
)

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type BotCommandScope struct {
	Type string `json:"type"`
}

type BotCommandScopeChat struct {
	Type   string      `json:"type"`
	ChatId interface{} `json:"chat_id"`
}

type CallbackQuery struct {
	Id              string  `json:"id"`
	From            User    `json:"from"`
	Message         Message `json:"message"`
	InlineMessageId string  `json:"inline_message_id"`
	ChatInstance    string  `json:"chat_instance"`
	Data            string  `json:"data"`
	GameShortName   string  `json:"game_short_name"`
}

func (cq CallbackQuery) Answer(ctx context.Context, b Bot, Text string) (MessageResponse, error) {
	acqr := AnswerCallbackQueryRequest{}
	acqr.CallbackQueryId = cq.Id
	acqr.Text = Text
	return b.Send(ctx, acqr)
}

type Chat struct {
	Id                    interface{}     `json:"id"`
	Type                  string          `json:"type"`
	Title                 string          `json:"title"`
	Username              string          `json:"username"`
	FirstName             string          `json:"first_name"`
	LastName              string          `json:"last_name"`
	Photo                 ChatPhoto       `json:"photo"`
	Bio                   string          `json:"bio"`
	Description           string          `json:"description"`
	InviteLink            string          `json:"invite_link"`
	PinnedMessage         *Message        `json:"pinned_message"`
	Permissions           ChatPermissions `json:"permissions"`
	SlowModeDelay         int             `json:"slow_mode_delay"`
	MessageAutoDeleteTime int             `json:"message_auto_delete_time"`
	StickerSetName        string          `json:"sticker_set_name"`
	CanSetStickerSet      bool            `json:"can_set_sticker_set"`
	LinkedChatId          int             `json:"linked_chat_id"`
	Location              ChatLocation    `json:"location"`
}

type ChatLocation struct {
	Location Location `json:"location"`
	Address  string   `json:"address"`
}

type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages"`
	CanSendMediaMessages  bool `json:"can_send_media_messages"`
	CanSendPolls          bool `json:"can_send_polls"`
	CanSendOtherMessages  bool `json:"can_send_other_messages"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews"`
	CanChangeInfo         bool `json:"can_change_info"`
	CanInviteUsers        bool `json:"can_invite_users"`
	CanPinMessages        bool `json:"can_pin_messages"`
}

type ChatPhoto struct {
	SmallFileId       string `json:"small_file_id"`
	SmallFileUniqueId string `json:"small_file_unique_id"`
	BigFileId         string `json:"big_file_id"`
	BigFileUniqueId   string `json:"big_file_unique_id"`
}

type ChatShared struct {
	RequestId int `json:"request_id"`
	ChatId    int `json:"chat_id"`
}

type KeyboardButton struct {
	Text        string                    `json:"text"`
	RequestChat KeyboardButtonRequestChat `json:"request_chat"`
}

type KeyboardButtonRequestChat struct {
	RequestId   int  `json:"request_id"`
	BotIsMember bool `json:"bot_is_member"`
}

type InlineKeyboardButton struct {
	Text                         string `json:"text"`
	Url                          string `json:"url,omitempty"`
	CallbackData                 string `json:"callback_data,omitempty"`
	SwitchInlineQuery            string `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat,omitempty"`
	Pay                          bool   `json:"pay,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type Location struct {
	Longitude            float32 `json:"longitude"`
	Latitude             float32 `json:"latitude"`
	HorizontalAccuracy   float32 `json:"horizontal_accuracy"`
	LivePeriod           int     `json:"live_period"`
	Heading              int     `json:"heading"`
	ProximityAlertRadius int     `json:"proximity_alert_radius"`
}

type Message struct {
	MessageId             int             `json:"message_id"`
	From                  User            `json:"from"`
	SenderChat            Chat            `json:"sender_chat"`
	Date                  int             `json:"date"`
	Chat                  Chat            `json:"chat"`
	ForwardFrom           User            `json:"forward_from"`
	ForwardFromChat       Chat            `json:"forward_from_chat"`
	ForwardFromMessage_id int             `json:"forward_from_message_id"`
	ForwardSignature      string          `json:"forward_signature"`
	ForwardSenderName     string          `json:"forward_sender_name"`
	ForwardDate           int             `json:"forward_date"`
	ReplyToMessage        *Message        `json:"reply_to_message"`
	ChatShared            ChatShared      `json:"chat_shared"`
	ViaBot                User            `json:"via_bot"`
	EditDate              int             `json:"edit_date"`
	MediaGroupId          int             `json:"media_group_id"`
	AuthorSignature       string          `json:"author_signature"`
	Text                  string          `json:"text"`
	Entities              []MessageEntity `json:"entities"`
	Caption               string          `json:"caption"`
	CaptionEntities       []MessageEntity `json:"caption_entities"`
	ReplyMarkup           interface{}     `json:"reply_markup"`
}

func (msg Message) GetCommand() string {
	re, _ := regexp.Compile(`^/([a-zA-Z0-9_]*)`)

	matches := re.FindStringSubmatch(msg.Text)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func (msg Message) IsCommand() bool {
	return msg.GetCommand() != ""
}

func (msg Message) DeleteMessage(ctx context.Context, b Bot) (MessageResponse, error) {
	return b.Send(ctx, DeleteMessageRequest{ChatId: msg.Chat.Id, MessageId: msg.MessageId})
}

// Edit
//
// Edit message with current Text and ReplyMarkup
func (msg Message) Edit(ctx context.Context, tb Bot) (MessageResponse, error) {
	return tb.Send(ctx,
		EditMessageTextRequest{ChatId: msg.Chat.Id, MessageId: msg.MessageId, Text: msg.Text, ReplyMarkup: msg.ReplyMarkup})
}

// EditKeyboard
//
// Edit message text
func (msg Message) EditKeyboard(ctx context.Context, tb Bot, kbd InlineKeyboardMarkup) (MessageResponse, error) {
	return tb.Send(ctx,
		EditMessageReplyMarkup{ChatId: msg.Chat.Id, MessageId: msg.MessageId, ReplyMarkup: kbd})
}

// EditMR
//
// Edit message with EditMessageTextRequest
func (msg Message) EditMR(ctx context.Context, tb Bot, emr EditMessageTextRequest) (MessageResponse, error) {
	emr.ChatId = msg.Chat.Id
	emr.MessageId = msg.MessageId
	return tb.Send(ctx, emr)
}

// EditText
//
// Edit message text
func (msg Message) EditText(ctx context.Context, tb Bot, text string) (MessageResponse, error) {
	return tb.Send(ctx,
		EditMessageTextRequest{ChatId: msg.Chat.Id, MessageId: msg.MessageId, Text: text})
}

func (msg Message) ReplyText(ctx context.Context, b Bot, text string) (MessageResponse, error) {
	return b.Send(ctx,
		MessageRequest{ChatId: msg.Chat.Id, ReplyToMessageId: msg.MessageId, Text: text})
}

func (msg Message) ReplyMR(ctx context.Context, b Bot, mr MessageRequest) (MessageResponse, error) {
	mr.ChatId = msg.Chat.Id
	mr.ReplyToMessageId = msg.MessageId
	return b.Send(ctx, mr)
}

// Send
//
// Send message with current Text and ReplyMarkup to same chat
func (msg Message) Send(ctx context.Context, b Bot) (MessageResponse, error) {
	return b.Send(ctx,
		MessageRequest{ChatId: msg.Chat.Id, Text: msg.Text, ReplyMarkup: msg.ReplyMarkup})
}

// SendText
//
// Send only text message to same chat
func (msg Message) SendText(ctx context.Context, b Bot, text string) (MessageResponse, error) {
	return b.Send(ctx,
		MessageRequest{ChatId: msg.Chat.Id, Text: text})
}

type MessageEntity struct {
	Type     string `json:"type"`
	Offset   int    `json:"offset"`
	Length   int    `json:"length"`
	Url      string `json:"url,omitempty"`
	User     *User  `json:"user,omitempty"`
	Language string `json:"language,omitempty"`
}
type ReplyKeyboardMarkup struct {
	Keyboard [][]KeyboardButton `json:"keyboard"`
}

type Update struct {
	UpdateId          int           `json:"update_id"`
	Message           Message       `json:"message"`
	EditedMessage     Message       `json:"edited_message"`
	ChannelPost       Message       `json:"channel_post"`
	EditedChannelPost Message       `json:"edited_channel_post"`
	CallbackQuery     CallbackQuery `json:"callback_query"`
}

type User struct {
	Id                      int    `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name,omitempty"`
	LastName                string `json:"last_name,omitempty"`
	UserName                string `json:"username,omitempty"`
	LanguageCode            string `json:"language_code,omitempty"`
	CanJoinGroups           bool   `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries,omitempty"`
}
