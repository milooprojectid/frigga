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
	var event Event
	if err := ctx.ReadJSON(&event); err != nil {
		return []Event{}, err
	}
	return []Event{
		event,
	}, nil
}

// EventReplier ...
func EventReplier(replyToken string, message string, accessToken string) error {
	payload := map[string]interface{}{
		"replyToken": replyToken,
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
	req.Header.Set("Authorization", "Bearer "+accessToken)

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
