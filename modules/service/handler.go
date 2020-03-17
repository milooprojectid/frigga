package service

import (
	"frigga/modules/line"
	messenger "frigga/modules/messenger"
	repo "frigga/modules/repository"
	telegram "frigga/modules/telegram"
	"strconv"
	"strings"

	"github.com/kataras/iris"
)

func replyWorker(subscription repo.Covid19Subcription, data map[string]int) {
	switch subscription.Provider {
	case "telegram":
		{
			bold := []string{}
			for key, val := range data {
				bold = append(bold, "<b>"+key+" "+strconv.Itoa(val)+"</b>")
			}
			messageType := "HTML"
			payload := telegram.MessageReply{
				Text:      strings.Join(bold, "\n"),
				ChatID:    subscription.Token,
				ParseMode: &messageType,
			}
			telegram.SendMessages(payload)
		}

	case "line":
		{
			bold := []string{}
			for key, val := range data {
				bold = append(bold, key+" "+strconv.Itoa(val))
			}
			payload := map[string]interface{}{
				"to": subscription.Token,
				"messages": []line.ReplyMessage{
					line.ReplyMessage{
						Text: strings.Join(bold, "\n"),
						Type: "text",
					},
				},
			}
			line.SendMessages(payload)
		}

	case "messenger":
		{
			bold := []string{}
			for key, val := range data {
				bold = append(bold, key+" "+strconv.Itoa(val))
			}
			payload := messenger.SendPayload{
				MessagingType: "RESPONSE",
				Recipient: messenger.SendPayloadRecipient{
					ID: subscription.Token,
				},
				Message: messenger.SendPayloadMessage{
					Text: strings.Join(bold, "\n"),
				},
			}
			messenger.SendMessages(payload)
		}
	}

}

// SendNotificationToSubscriptionHandler ...
func SendNotificationToSubscriptionHandler(ctx iris.Context) {
	covid19data, _ := repo.GetCovid19Data()
	covid19subcriptions, _ := repo.GetCovid19Subscription()

	for _, sub := range covid19subcriptions {
		go replyWorker(sub, covid19data)
	}

	ctx.JSON(map[string]interface{}{
		"message": "subscription broadcasted",
		"data":    len(covid19subcriptions),
	})
}
