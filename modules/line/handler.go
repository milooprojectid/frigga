package line

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
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
func EventReplier(messages []string, quickReply *QuickReply, replyToken string) error {
	var replyMessages []ReplyMessage
	accessToken := os.Getenv("LINE_TOKEN")

	for _, message := range messages {
		replyMessages = append(replyMessages, ReplyMessage{
			Type:       "text",
			Text:       message,
			QuickReply: quickReply,
		})
	}

	payload := map[string]interface{}{
		"to":       replyToken,
		"messages": replyMessages,
	}
	requestBody, _ := json.Marshal(payload)

	client := &http.Client{Timeout: time.Second * 60}
	req, err := http.NewRequest("POST", apiURL+"/message/push", bytes.NewBuffer(requestBody))
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

// GetUserName ...
func GetUserName(userID string) (string, error) {
	var profile userProfile
	token := os.Getenv("LINE_TOKEN")

	client := &http.Client{Timeout: time.Second * 60}
	req, err := http.NewRequest("GET", apiURL+"/profile/"+userID, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &profile)
	defer response.Body.Close()

	return profile.DisplayName, nil
}

// GetCommandQuickReply ...
func GetCommandQuickReply() *QuickReply {
	quick := &QuickReply{
		Items: []QuickReplyItem{
			QuickReplyItem{
				Type: "action",
				Action: ItemAction{
					Type:  "message",
					Label: "sentiment",
					Text:  "/sentiment",
				},
			},
			QuickReplyItem{
				Type: "action",
				Action: ItemAction{
					Type:  "message",
					Label: "summarize",
					Text:  "/summarize",
				},
			},
			QuickReplyItem{
				Type: "action",
				Action: ItemAction{
					Type:  "message",
					Label: "covid-19 status",
					Text:  "/corona",
				},
			},
		},
	}
	return quick
}
