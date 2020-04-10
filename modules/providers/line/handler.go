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
			payload = MediaReplyMessage{
				Type:               "audio",
				OriginalContentURL: message.Text,
				QuickReply:         quickReply,
			}
		case c.VideoMessageType:
			payload = MediaReplyMessage{
				Type:               "video",
				OriginalContentURL: message.Text,
				QuickReply:         quickReply,
			}
		case c.ImageMessageType:
			payload = MediaReplyMessage{
				Type:               "image",
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
			QuickReplyItem{
				Type: "action",
				Action: ItemAction{
					Type:  "message",
					Label: "covid-19 notif",
					Text:  "/subscov19",
				},
			},
		},
	}
	return quick
}
