package messenger

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/kataras/iris"
)

const apiURL = "https://graph.facebook.com/v6.0"

// EventAdapter ...
func EventAdapter(ctx iris.Context) ([]Entry, error) {
	var event Event
	if err := ctx.ReadJSON(&event); err != nil {
		return event.Entry, err
	}
	return event.Entry, nil
}

// EventReplier ...
func EventReplier(PSID string, message string, token string) error {
	payload := SendPayload{
		MessagingType: "RESPONSE",
		Recipient: SendPayloadRecipient{
			ID: PSID,
		},
		Message: SendPayloadMessage{
			Text: message,
		},
	}

	requestBody, _ := json.Marshal(payload)
	if _, err := http.Post(apiURL+"/me/messages?access_token="+token, "application/json", bytes.NewBuffer(requestBody)); err != nil {
		return err
	}
	return nil
}
