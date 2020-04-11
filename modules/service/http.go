package service

import (
	c "frigga/modules/common"
	"frigga/modules/providers/line"
	messenger "frigga/modules/providers/messenger"
	telegram "frigga/modules/providers/telegram"
	repo "frigga/modules/repository"
	"frigga/modules/template"
	"time"

	"github.com/kataras/iris"
	"gopkg.in/go-playground/validator.v9"
)

func replyWorker(message c.Message, provider string, replyToken string) {
	switch provider {
	case telegram.Name:
		telegram.EventReplier(message, nil, replyToken)

	case line.Name:
		{
			quickReply := line.GetCommandQuickReply()
			line.EventReplier([]c.Message{
				message,
			}, quickReply, replyToken)
		}

	case messenger.Name:
		messenger.EventReplier(message, nil, replyToken)
	}

}

// SendNotificationToSubscriptionHandler ...
func SendNotificationToSubscriptionHandler(ctx iris.Context) {
	covid19data, _ := repo.GetCovid19Data()
	covid19subcriptions, _ := repo.GetCovid19Subscription()

	templatePayload := map[string]interface{}{
		"Date":      time.Now().Format("02 Jan 2006"),
		"Confirmed": covid19data.Confirmed,
		"Recovered": covid19data.Recovered,
		"Deceased":  covid19data.Deceased,
		"Treated":   covid19data.Confirmed - covid19data.Recovered - covid19data.Deceased,
	}
	message := template.ProcessFile("storage/covid19.tmpl", templatePayload)

	for _, sub := range covid19subcriptions {
		messagePayload := c.GenerateTextMessage(message)
		go replyWorker(messagePayload, sub.Provider, sub.Token)
	}

	ctx.JSON(map[string]interface{}{
		"message": "subscription broadcasted",
		"data":    len(covid19subcriptions),
	})
}

// SendBroadcastMessageHandler ...
func SendBroadcastMessageHandler(ctx iris.Context) {
	var data BroadcastMessage
	ctx.ReadJSON(&data)

	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		ctx.JSON(map[string]interface{}{
			"message": "validation fail",
			"data":    err.Error(),
		})
		return
	}

	sessions, _ := repo.GetAllBotSessions(data.Mode == "prod")
	sessionLength := len(sessions)
	start := 0
	batch := 50

	var message c.Message
	switch data.Type {
	case c.ImageMessageType:
		message = c.GenerateImageMessage(data.Body)
	default:
		message = c.GenerateTextMessage(data.Body)
	}

	// loop per batch
	for {
		currentStart := start
		end := start + batch
		if end > sessionLength {
			end = sessionLength
		}

		for i := currentStart; i < end; i++ {
			go replyWorker(message, sessions[i].Provider, sessions[i].UserID)
		}

		if end == sessionLength {
			break
		} else {
			start = end
			time.Sleep(2 * time.Second)
		}
	}

	ctx.JSON(map[string]interface{}{
		"message": "message broadcasted",
		"data":    sessionLength,
	})
	return
}
