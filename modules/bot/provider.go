package bot

import (
	"os"
	"strconv"

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
	}

	return provider
}
