package messenger

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/kataras/iris"
)

const apiURL = "https://graph.facebook.com/v6.0"

// EventAdapter ...
func EventAdapter(ctx iris.Context) ([]Messaging, error) {
	var event Event
	var messagings []Messaging
	if err := ctx.ReadJSON(&event); err != nil {
		return messagings, err
	}

	for _, en := range event.Entry {
		for _, messaging := range en.Messaging {
			messagings = append(messagings, messaging)
		}
	}

	return messagings, nil
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
