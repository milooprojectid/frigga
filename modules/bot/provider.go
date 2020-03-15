package bot

import (
	"log"
	"os"
	"strconv"
	"strings"

	line "frigga/modules/line"
	messenger "frigga/modules/messenger"
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

func trim(input string) string {
	return strings.TrimRight(input, "\t \n")
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
							Message: trim(update.Message.Text),
							Token:   strconv.Itoa(update.Message.Chat.ID),
						})
					}
					return events, nil
				},
				EventReplier: func(rep eventReply) {
					var replyMarkup *telegram.ReplyMarkup
					// if !rep.isTrigger() {
					// 	replyMarkup = telegram.GetCommandQuickReply()
					// }
					for _, message := range rep.Messages {
						telegram.EventReplier(message, replyMarkup, rep.Token, token)
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
								Message: trim(update.Message.Text),
								Token:   update.Source.UserID,
							})
						} else if update.Type == "follow" {
							events = append(events, Event{
								ID:      update.Source.UserID,
								Message: "/start",
								Token:   update.Source.UserID,
							})
						}
					}
					return events, nil
				},
				EventReplier: func(rep eventReply) {
					var quickReply *line.QuickReply
					if !rep.isTrigger() {
						quickReply = line.GetCommandQuickReply()
					}

					err := line.EventReplier(rep.Messages, quickReply, rep.Token, token)
					if err != nil {
						log.Printf(err.Error())
					}
				},
			}

		}
	case "messenger":
		{
			token := os.Getenv("MESSENGER_TOKEN")
			provider = Provider{
				Name:        name,
				AccessToken: token,
				EventAdapter: func(ctx iris.Context) ([]Event, error) {
					var events []Event
					messages, _ := messenger.EventAdapter(ctx)
					for _, message := range messages {
						events = append(events, Event{
							ID:      message.Sender.ID,
							Message: trim(message.Message.Text),
							Token:   message.Sender.ID,
						})
					}
					return events, nil
				},
				EventReplier: func(rep eventReply) {
					var quickReplies *[]messenger.QuickReply
					if !rep.isTrigger() {
						quickReplies = messenger.GetCommandQuickReply()
					}

					for _, message := range rep.Messages {
						messenger.EventReplier(message, quickReplies, rep.Token, token)
					}
				},
			}
		}
	}

	return provider
}
