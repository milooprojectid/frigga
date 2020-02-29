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
			provider = Provider{
				Name:        name,
				AccessToken: os.Getenv("TELEGRAM_TOKEN"),
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
					telegram.EventReplier(rep.Token, rep.Message, os.Getenv("TELEGRAM_TOKEN"))
				},
			}

		}

	case "line":
		{
			provider = Provider{
				Name:        name,
				AccessToken: os.Getenv("LINE_TOKEN"),
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
					line.EventReplier(rep.Message, rep.Token)
				},
			}

		}
	}

	return provider
}
