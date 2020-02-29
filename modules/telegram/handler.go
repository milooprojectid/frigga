package telegram

import (
	"bytes"
	"encoding/json"
	"net/http"

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
func EventReplier(chatID string, message string, token string) error {
	payload := map[string]string{"chat_id": chatID, "text": message}
	requestBody, _ := json.Marshal(payload)
	if _, err := http.Post(apiURL+token+"/sendMessage", "application/json", bytes.NewBuffer(requestBody)); err != nil {
		return err
	}
	return nil
}
