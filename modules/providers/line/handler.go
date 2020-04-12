package line

import (
	"bytes"
	"encoding/json"
	c "frigga/modules/common"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
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

// SendMessages ...
func SendMessages(payload interface{}) error {
	token := os.Getenv("LINE_TOKEN")
	requestBody, _ := json.Marshal(payload)

	client := &http.Client{Timeout: time.Second * 60}
	req, err := http.NewRequest("POST", apiURL+"/message/push", bytes.NewBuffer(requestBody))
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

// EventReplier ...
func EventReplier(messages []c.Message, quickReply *QuickReply, replyToken string) error {
	var replyMessages []interface{}

	for _, message := range messages {
		var payload interface{}

		switch message.Type {
		case c.AudioMessageType:
			payload = AudioReplyMessage{
				Type:               "audio",
				OriginalContentURL: message.Text,
				QuickReply:         quickReply,
				Duration:           60000,
			}
		case c.VideoMessageType:
			payload = VideoReplyMessage{
				Type:               "video",
				OriginalContentURL: message.Text,
				PreviewImageURL:    "https://miloo.id/assets/img/thumbnail.png",
				QuickReply:         quickReply,
			}
		case c.ImageMessageType:
			payload = ImageReplyMessage{
				Type:               "image",
				PreviewImageURL:    "https://miloo.id/assets/img/thumbnail.png",
				OriginalContentURL: message.Text,
				QuickReply:         quickReply,
			}
		case c.LocationMessageType:
			{
				splitted := strings.Split(message.Text, ",")
				lat, _ := strconv.ParseFloat(splitted[0], 64)
				lon, _ := strconv.ParseFloat(splitted[1], 64)
				payload = LocationReplyMessage{
					Type:       "location",
					Title:      "Location",
					Address:    message.Text,
					Latitude:   lat,
					Longitude:  lon,
					QuickReply: quickReply,
				}
			}
		default:
			payload = TextReplyMessage{
				Type:       "text",
				Text:       message.Text,
				QuickReply: quickReply,
			}
		}

		replyMessages = append(replyMessages, payload)
	}

	payloadData := map[string]interface{}{
		"to":       replyToken,
		"messages": replyMessages,
	}

	return SendMessages(payloadData)
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
func GetCommandQuickReply(commands []string) *QuickReply {
	var items []QuickReplyItem

	for _, path := range commands {
		items = append(items, QuickReplyItem{
			Type: "action",
			Action: ItemAction{
				Type:  "message",
				Label: path,
				Text:  path,
			},
		})
	}

	return &QuickReply{
		Items: items,
	}
}
