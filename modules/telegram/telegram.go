package telegram

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/kataras/iris"
)

const apiURL = "https://api.telegram.org/bot"

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

type messageReply struct {
	ChatID  int
	Message string
}

// Process ...
func (e *Event) Process(c chan messageReply) {
	c <- messageReply{
		ChatID:  e.Message.Chat.ID,
		Message: e.Message.Text,
	}
}

// Handler ...
func (b *Bot) Handler(ctx iris.Context) {
	var events []Event
	if err := ctx.ReadJSON(&events); err != nil {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.WriteString(err.Error())
		return
	}

	replies := make(chan messageReply)

	for _, event := range events {
		go event.Process(replies)
	}

	for i := 0; i < len(events); i++ {
		r := <-replies
		go b.Reply(r.ChatID, r.Message)
	}

	ctx.JSON(events)
	return
}

// Reply ...
func (b *Bot) Reply(chatID int, message string) error {
	payload := map[string]interface{}{"chat_id": chatID, "text": message}
	requestBody, _ := json.Marshal(payload)

	if _, err := http.Post(apiURL+b.Token+"/sendMessage", "application/json", bytes.NewBuffer(requestBody)); err != nil {
		return err
	}
	return nil
}

// NewBot returns new bot instance
func NewBot(token string) Bot {
	return Bot{
		Token: token,
	}
}
