package bot

import (
	"log"
	"os"
	"strconv"
	"strings"

	c "frigga/modules/common"
	"frigga/modules/providers/discord"
	line "frigga/modules/providers/line"
	messenger "frigga/modules/providers/messenger"
	telegram "frigga/modules/providers/telegram"

	"github.com/kataras/iris"
)

// Provider ...
type Provider struct {
	Name           string
	AccessToken    string
	EventAdapter   func(ctx iris.Context) ([]BotEvent, error)
	EventReplier   func(rep BotReply)
	StaticResponse map[string]interface{}
}

func trim(input string) string {
	return strings.TrimRight(input, "\t \n")
}

func getDefaultResponse() map[string]interface{} {
	return map[string]interface{}{
		"message": "event dispatched",
	}
}

// GetProvider ...
func GetProvider(name string) Provider {
	var provider Provider

	switch name {
	case telegram.Name:
		{
			provider = Provider{
				Name:        name,
				AccessToken: os.Getenv("TELEGRAM_TOKEN"),
				EventAdapter: func(ctx iris.Context) ([]BotEvent, error) {
					var events []BotEvent
					updates, _ := telegram.EventAdapter(ctx)
					for _, update := range updates {
						message := telegram.CheckIncomingMessage(update)
						if message != "" {
							events = append(events, BotEvent{
								ID:       strconv.Itoa(update.Message.Chat.ID),
								Message:  c.GenerateTextMessage(trim(message)),
								Token:    strconv.Itoa(update.Message.Chat.ID),
								Provider: telegram.Name,
							})
						}
					}
					return events, nil
				},
				EventReplier: func(rep BotReply) {
					var replyMarkup *telegram.ReplyMarkup
					for _, message := range rep.Messages {
						telegram.EventReplier(message, replyMarkup, rep.Token)
					}
				},
				StaticResponse: getDefaultResponse(),
			}

		}

	case line.Name:
		{
			provider = Provider{
				Name:        name,
				AccessToken: os.Getenv("LINE_TOKEN"),
				EventAdapter: func(ctx iris.Context) ([]BotEvent, error) {
					var events []BotEvent
					updates, _ := line.EventAdapter(ctx)
					for _, update := range updates {
						if update.Type == "message" {
							events = append(events, BotEvent{
								ID:       update.Source.UserID,
								Message:  c.GenerateTextMessage(trim(update.Message.Text)),
								Token:    update.Source.UserID,
								Provider: line.Name,
							})
						} else if update.Type == "follow" {
							events = append(events, BotEvent{
								ID:       update.Source.UserID,
								Message:  c.GenerateTextMessage("/start"),
								Token:    update.Source.UserID,
								Provider: line.Name,
							})
						}
					}
					return events, nil
				},
				EventReplier: func(rep BotReply) {
					var quickReply *line.QuickReply
					err := line.EventReplier(rep.Messages, quickReply, rep.Token)
					if err != nil {
						log.Printf(err.Error())
					}
				},
				StaticResponse: getDefaultResponse(),
			}

		}
	case messenger.Name:
		{
			provider = Provider{
				Name:        name,
				AccessToken: os.Getenv("MESSENGER_TOKEN"),
				EventAdapter: func(ctx iris.Context) ([]BotEvent, error) {
					var events []BotEvent
					messages, _ := messenger.EventAdapter(ctx)
					for _, message := range messages {
						events = append(events, BotEvent{
							ID:       message.Sender.ID,
							Message:  c.GenerateTextMessage(trim(message.Message.Text)),
							Token:    message.Sender.ID,
							Provider: messenger.Name,
						})
					}
					return events, nil
				},
				EventReplier: func(rep BotReply) {
					var quickReplies *[]messenger.QuickReply
					for _, message := range rep.Messages {
						messenger.EventReplier(message, quickReplies, rep.Token)
					}
				},
				StaticResponse: getDefaultResponse(),
			}
		}

	case discord.Name:
		{
			provider = Provider{
				Name:        name,
				AccessToken: os.Getenv("DISCORD_PUBLIC_KEY"),
				EventAdapter: func(ctx iris.Context) ([]BotEvent, error) {
					var events []BotEvent
					requests, _ := discord.EventAdapter(ctx)
					for _, request := range requests {
						events = append(events, BotEvent{
							ID:       request.GetId(),
							Message:  c.GenerateTextMessage(trim(request.Data.GetInlineCommand())),
							Token:    request.Token,
							Provider: discord.Name,
						})
					}
					return events, nil
				},
				EventReplier: func(rep BotReply) {
					for i, message := range rep.Messages {
						discord.EventReplier(message, rep.Token, i == 0)
					}
				},
				StaticResponse: map[string]interface{}{
					"type": 5, // deffered
					"data": map[string]interface{}{
						"content": "Processing command!",
					},
				},
			}
		}

	}

	return provider
}
