package telegram

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	repo "frigga/modules/repository"
	s "frigga/modules/service"

	"github.com/kataras/iris"
)

const apiURL = "https://api.telegram.org/bot"
const commands = "You can control me by sending these commands:\n/sentiment - run sentiment analysis on a text\n/summarize - summarise a content of a link or text\n/cancel - terminate currently running command"

type sentimentResult struct {
	Data struct {
		Class       int    `json:"class"`
		Description string `json:"description"`
	} `json:"data"`
	Message string `json:"message"`
}

type summarizationResult struct {
	Data struct {
		Summary string `json:"summary"`
	} `json:"data"`
	Message string `json:"message"`
}

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
	var message string
	var input string = e.Message.Text
	var chatID string = strconv.Itoa(e.Message.Chat.ID)

	switch input {
	case "/start":
		{
			message = "Hi im Miloo\n" + commands
			repo.InitSession(chatID, e.Message.Chat.FirstName, e.Message.Chat.LastName)
		}
	case "/sentiment":
		{
			message = "Type the statement you want to analize"
			repo.UpdateSession(chatID, input)
		}
	case "/summarize":
		{
			message = "Type the statement or url you want to summarise"
			repo.UpdateSession(chatID, input)
		}
	case "/cancel":
		{
			if command, _ := repo.GetSession(chatID); command == "" {
				message = "No active command to cancel. I wasn't doing anything anyway. Zzzzz..."
			} else if command, _ := repo.GetSession(chatID); command == "" {
				message = "Command cancelled"
				repo.UpdateSession(chatID, "")
			}
		}
	default:
		{
			message = ""
		}
	}

	if message == "" {
		if len(input) == 0 {
			message = "You have to type something ._."
		} else if string(input[0]) == "/" {
			message = "I cant understand that command\n" + commands
		}

		if command, _ := repo.GetSession(chatID); command == "" {
			message = "Hmmm, theres no active command >_<"
		} else if command == "/summarize" {
			var result summarizationResult
			service, _ := s.GetService("storm")
			service.Call("summarizeText", map[string]string{"text": input}, &result)
			message = result.Data.Summary
		} else if command == "/sentiment" {
			var result sentimentResult
			service, _ := s.GetService("morbius")
			service.Call("sentiment", map[string]string{"text": input}, &result)
			message = result.Data.Description
		}

		repo.UpdateSession(chatID, "")
	}

	c <- messageReply{
		ChatID:  e.Message.Chat.ID,
		Message: message,
	}
}

// Handler ...
func (b *Bot) Handler(ctx iris.Context) {
	var event Event
	if err := ctx.ReadJSON(&event); err != nil {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.WriteString(err.Error())
		return
	}

	replies := make(chan messageReply)
	go event.Process(replies)

	r := <-replies
	go b.Reply(r.ChatID, r.Message)

	ctx.JSON(event)
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
