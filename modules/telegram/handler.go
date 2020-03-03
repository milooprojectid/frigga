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
func EventReplier(chatID string, message string, token string) error {
	payload := map[string]string{"chat_id": chatID, "text": message}
	requestBody, _ := json.Marshal(payload)
	if _, err := http.Post(apiURL+token+"/sendMessage", "application/json", bytes.NewBuffer(requestBody)); err != nil {
		return err
	}
	return nil
}

// GetUserName ...
func GetUserName(chatID string) (string, error) {
	var chat chat
	var name string

	token := os.Getenv("TELEGRAM_TOKEN")
	payload := map[string]string{"chat_id": chatID}
	requestBody, _ := json.Marshal(payload)

	response, err := http.Post(apiURL+token+"/getChat", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &chat)
	defer response.Body.Close()

	name = chat.FirstName
	if chat.LastName != "" {
		name = name + "" + chat.LastName
	}

	return name, nil
}
