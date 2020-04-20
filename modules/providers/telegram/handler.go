package telegram

import (
	"bytes"
	"encoding/json"
	c "frigga/modules/common"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

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

// SendMessages ...
func SendMessages(payload interface{}, messageType string) error {
	token := os.Getenv("TELEGRAM_TOKEN")
	requestBody, _ := json.Marshal(payload)
	var url string

	switch messageType {
	case c.TextMessageType:
		url = "sendMessage"

	case c.AudioMessageType:
		url = "sendAudio"

	case c.VideoMessageType:
		url = "sendVideo"

	case c.ImageMessageType:
		url = "sendPhoto"

	case c.LocationMessageType:
		url = "sendLocation"
	}

	if _, err := http.Post(apiURL+token+"/"+url, "application/json", bytes.NewBuffer(requestBody)); err != nil {
		return err
	}

	return nil
}

// EventReplier ...
func EventReplier(message c.Message, replyMarkup *ReplyMarkup, chatID string) error {
	var payload interface{}

	switch message.Type {
	case c.AudioMessageType:
		payload = AudioMessageReply{
			Audio:       message.Text,
			ChatID:      chatID,
			ReplyMarkup: replyMarkup,
		}
	case c.VideoMessageType:
		payload = VideoMessageReply{
			Video:       message.Text,
			ChatID:      chatID,
			ReplyMarkup: replyMarkup,
		}
	case c.ImageMessageType:
		payload = ImageMessageReply{
			Photo:                 message.Text,
			ChatID:                chatID,
			DisableWebPagePreview: true,
			ReplyMarkup:           replyMarkup,
		}
	case c.LocationMessageType:
		{
			splitted := strings.Split(message.Text, ",")
			lat, _ := strconv.ParseFloat(splitted[0], 64)
			lon, _ := strconv.ParseFloat(splitted[1], 64)
			payload = LocationMessageReply{
				ChatID:      chatID,
				Latitude:    lat,
				Longitude:   lon,
				ReplyMarkup: replyMarkup,
			}
		}
	default:
		payload = TextMessageReply{
			Text:                  message.Text,
			ChatID:                chatID,
			DisableWebPagePreview: true,
			ReplyMarkup:           replyMarkup,
		}
	}

	return SendMessages(payload, message.Type)
}

// GetUserName ...
func GetUserName(chatID string) (string, error) {
	var profile struct {
		OK     bool `json:"ok"`
		Result chat `json:"result"`
	}
	var name string

	token := os.Getenv("TELEGRAM_TOKEN")
	payload := map[string]string{"chat_id": chatID}
	requestBody, _ := json.Marshal(payload)

	response, err := http.Post(apiURL+token+"/getChat", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &profile)
	defer response.Body.Close()

	name = profile.Result.FirstName
	if profile.Result.LastName != "" {
		name = name + " " + profile.Result.LastName
	}

	if name == "" {
		name = profile.Result.Title
	}

	return name, nil
}

// CheckIncomingMessage ...
func CheckIncomingMessage(update Update) string {
	message := update.Message.Text

	// check if not a reply
	if update.Message.ReplyToMessage != nil {
		return ""
	}

	// check if containt bot identifier
	splitted := strings.SplitN(message, "@", 2)
	if len(splitted) != 2 {
		return message
	}

	// check if bot identifier id correct
	if splitted[1] != "miloo_bot" {
		return ""
	}

	return splitted[0]
}

// GetCommandQuickReply ...
func GetCommandQuickReply(commands []string) *ReplyMarkup {
	var quicks []InlineKeyboard

	for _, path := range commands {
		quicks = append(quicks, InlineKeyboard{
			Text:         path,
			CallbackData: path,
		})
	}

	return &ReplyMarkup{
		InlineKeyboard: [][]InlineKeyboard{
			quicks,
		},
	}
}
