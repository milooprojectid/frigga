package telegram

// Bot ...
type Bot struct {
	Token string
}

type user struct {
	ID        int    `json:"id"`
	IsBOT     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

type chat struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type message struct {
	MessageID int    `json:"message_id"`
	From      user   `json:"from"`
	Date      int    `json:"date"`
	Chat      chat   `json:"chat"`
	Text      string `json:"text,omitempty"`
}

// Update ...
type Update struct {
	UpdateID int     `json:"update_id"`
	Message  message `json:"message"`
}

// InlineKeyboard ...
type InlineKeyboard struct {
	Text string  `json:"text"`
	URL  *string `json:"url,omitempty"`
}

// ReplyMarkup ...
type ReplyMarkup struct {
	InlineKeyboard []InlineKeyboard `json:"inline_keyboard"`
}

// MessageReply ...
type MessageReply struct {
	ChatID                string       `json:"chat_id"`
	Text                  string       `json:"text"`
	DisableWebPagePreview bool         `json:"disable_web_page_preview"`
	ParseMode             *string      `json:"parse_mode,omitempty"`
	ReplyMarkup           *ReplyMarkup `json:"reply_markup,omitempty"`
}
