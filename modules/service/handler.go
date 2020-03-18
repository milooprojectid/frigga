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
	bold := []string{}
	for key, val := range data {
		bold = append(bold, key+" "+strconv.Itoa(val))
	}
	message := "Covid-19 Report\n\n" + strings.Join(bold, "\n")

	switch subscription.Provider {
	case "telegram":
		{
			messageType := "HTML"
			payload := telegram.MessageReply{
				Text:      message,
				ChatID:    subscription.Token,
				ParseMode: &messageType,
			}
			telegram.SendMessages(payload)
		}

	case "line":
		{
			payload := map[string]interface{}{
				"to": subscription.Token,
				"messages": []line.ReplyMessage{
					line.ReplyMessage{
						Text: message,
						Type: "text",
					},
				},
			}
			line.SendMessages(payload)
		}

	case "messenger":
		{
			payload := messenger.SendPayload{
				MessagingType: "RESPONSE",
				Recipient: messenger.SendPayloadRecipient{
					ID: subscription.Token,
				},
				Message: messenger.SendPayloadMessage{
					Text: message,
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
