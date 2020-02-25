package telegram

import (
	"github.com/kataras/iris"
)

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

// Event ...
type Event struct {
	UpdateID int     `json:"update_id"`
	Message  message `json:"message"`
}

// Handler ...
func (b *Bot) Handler(ctx iris.Context) {
	var events []Event
	ctx.ReadJSON(&events)
	ctx.JSON(events)
}

// NewBot returns new bot instance
func NewBot(token string) Bot {
	return Bot{
		Token: token,
	}
}
