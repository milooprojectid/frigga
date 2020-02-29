package line

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/kataras/iris"
)

const apiURL = "https://api.line.me/v2/bot"

// EventAdapter ...
func EventAdapter(ctx iris.Context) ([]Event, error) {
	var events []Event
	if err := ctx.ReadJSON(&events); err != nil {
		return events, err
	}
	return events, nil
}

// EventReplier ...
func EventReplier(message string, token string) error {
	payload := map[string]interface{}{
		"replyToken": token,
		"messages": []map[string]string{
			{"type": "text", "text": message},
		},
	}
	requestBody, _ := json.Marshal(payload)

	client := &http.Client{Timeout: time.Second * 60}
	req, err := http.NewRequest("POST", apiURL+"/message/reply", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
