package telegram

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/kataras/iris"
)

const apiURL = "https://api.telegram.org/bot"

// EventAdapter ...
func EventAdapter(ctx iris.Context) ([]Update, error) {
	var update Update
	if err := ctx.ReadJSON(&update); err != nil {
		return []Update{
			update,
		}, err
	}
	return []Update{
		update,
	}, nil
}

// EventReplier ...
func EventReplier(message string, replyMarkup *ReplyMarkup, chatID string) error {
	token := os.Getenv("TELEGRAM_TOKEN")
	payload := MessageReply{
		Text:                  message,
		ChatID:                chatID,
		DisableWebPagePreview: true,
		ReplyMarkup:           replyMarkup,
	}

	requestBody, _ := json.Marshal(payload)
	if _, err := http.Post(apiURL+token+"/sendMessage", "application/json", bytes.NewBuffer(requestBody)); err != nil {
		return err
	}
	return nil
}

// GetUserName ...
func GetUserName(chatID string) (string, error) {
	var profile struct {
		OK     bool `json:"ok"`
		Result chat `json:"result"`
	}
	var name string

	token := os.Getenv("TELEGRAM_TOKEN")
	payload := map[string]string{"chat_id": chatID}
	requestBody, _ := json.Marshal(payload)

	response, err := http.Post(apiURL+token+"/getChat", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &profile)
	defer response.Body.Close()

	name = profile.Result.FirstName
	if profile.Result.LastName != "" {
		name = name + " " + profile.Result.LastName
	}

	return name, nil
}

// GetCommandQuickReply ...
func GetCommandQuickReply() *ReplyMarkup {
	quick := &ReplyMarkup{
		InlineKeyboard: []InlineKeyboard{
			InlineKeyboard{
				Text: "/summarize",
			},
			InlineKeyboard{
				Text: "/sentiment",
			},
			InlineKeyboard{
				Text: "/corona",
			},
		},
	}

	return quick
}
