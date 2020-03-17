package service

import (
	line "frigga/modules/line"
	repo "frigga/modules/repository"
	telegram "frigga/modules/telegram"

	"github.com/kataras/iris"
)

func replyWorker(subscription repo.Covid19Subcription, data map[string]int) {
	switch subscription.Provider {
	case "telegram":
		{
			line.EventReplier()
		}

	case "line":
		{
			telegram.EventReplier()
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
