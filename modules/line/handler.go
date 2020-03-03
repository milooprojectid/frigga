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
	var request Request
	if err := ctx.ReadJSON(&request); err != nil {
		return request.Events, err
	}
	return request.Events, nil
}

// EventReplier ...
func EventReplier(replyToken string, messages []string, accessToken string) error {
	var replyMessages []map[string]string
	for _, message := range messages {
		replyMessages = append(replyMessages, map[string]string{"type": "text", "text": message})
	}

	payload := map[string]interface{}{
		"replyToken": replyToken,
		"messages":   replyMessages,
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
