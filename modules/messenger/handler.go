package messenger

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/kataras/iris"
)

const apiURL = "https://graph.facebook.com"

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
func EventReplier(message string, quickReplies *[]QuickReply, PSID string, token string) error {
	payload := SendPayload{
		MessagingType: "RESPONSE",
		Recipient: SendPayloadRecipient{
			ID: PSID,
		},
		Message: SendPayloadMessage{
			Text:       message,
			QuickReply: quickReplies,
		},
	}

	requestBody, _ := json.Marshal(payload)
	if _, err := http.Post(apiURL+"/v6.0/me/messages?access_token="+token, "application/json", bytes.NewBuffer(requestBody)); err != nil {
		return err
	}
	return nil
}

// GetUserName ...
func GetUserName(PSID string) (string, error) {
	var profile userProfile
	var name string
	token := os.Getenv("MESSENGER_TOKEN")

	client := &http.Client{Timeout: time.Second * 60}
	req, err := http.NewRequest("GET", apiURL+"/"+PSID+"?access_token="+token+"&fields=first_name,last_name", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &profile)
	defer response.Body.Close()

	name = profile.FirstName
	if profile.LastName != "" {
		name = name + " " + profile.LastName
	}

	return name, nil
}

// GetCommandQuickReply ...
func GetCommandQuickReply() *[]QuickReply {
	quick := &[]QuickReply{
		QuickReply{
			ContentType: "text",
			Title:       "/sentiment",
			Payload:     "/sentiment",
		},
		QuickReply{
			ContentType: "text",
			Title:       "/summarize",
			Payload:     "/summarize",
		},
	}

	return quick
}