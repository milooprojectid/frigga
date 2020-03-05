package bot

import (
	"os"
	"strconv"

	"frigga/modules/line"
	telegram "frigga/modules/telegram"

	"github.com/kataras/iris"
)

// Provider ...
type Provider struct {
	Name         string
	AccessToken  string
	EventAdapter func(ctx iris.Context) ([]Event, error)
	EventReplier func(rep eventReply)
}

// GetProvider ...
func GetProvider(name string) Provider {
	var provider Provider

	switch name {
	case "telegram":
		{
			token := os.Getenv("TELEGRAM_TOKEN")
			provider = Provider{
				Name:        name,
				AccessToken: token,
				EventAdapter: func(ctx iris.Context) ([]Event, error) {
					var events []Event
					updates, _ := telegram.EventAdapter(ctx)
					for _, update := range updates {
						events = append(events, Event{
							ID:      strconv.Itoa(update.Message.Chat.ID),
							Message: update.Message.Text,
							Token:   strconv.Itoa(update.Message.Chat.ID),
						})
					}
					return events, nil
				},
				EventReplier: func(rep eventReply) {
					for _, message := range rep.Messages {
						telegram.EventReplier(rep.Token, message, token)
					}
				},
			}

		}

	case "line":
		{
			token := os.Getenv("LINE_TOKEN")
			provider = Provider{
				Name:        name,
				AccessToken: token,
				EventAdapter: func(ctx iris.Context) ([]Event, error) {
					var events []Event
					updates, _ := line.EventAdapter(ctx)
					for _, update := range updates {
						if update.Type == "message" {
							events = append(events, Event{
								ID:      update.Source.UserID,
								Message: update.Message.Text,
								Token:   update.ReplyToken,
							})
						} else if update.Type == "follow" {
							events = append(events, Event{
								ID:      update.Source.UserID,
								Message: "/start",
								Token:   update.ReplyToken,
							})
						}
					}
					return events, nil
				},
				EventReplier: func(rep eventReply) {
					var quickReply line.QuickReply
					if !rep.isTrigger() {
						quickReply = line.GetCommandQuickReply()
					}
					line.EventReplier(rep.Token, rep.Messages, &quickReply, token)
				},
			}

		}
	}

	return provider
}
