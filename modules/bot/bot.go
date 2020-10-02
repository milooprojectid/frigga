package bot

import (
	"github.com/kataras/iris"
)

// Bot ...
type Bot struct {
	Provider Provider
	Commands *commands
}

func (b *Bot) dispatch(e BotEvent, r chan BotReply) {
	if ok := e.IsInline(); ok {
		r <- Commands.ExecuteInline(e)
	} else {
		r <- Commands.Execute(e)
	}
}

// Handler will intercept incoming request and pass it to provider
func (b *Bot) Handler(ctx iris.Context) {
	events, _ := b.Provider.EventAdapter(ctx)
	isWithWorker := false

	if isWithWorker {
		for _, event := range events {
			DispatchBotEventWorker(event)
		}
	} else {
		replyChannel := make(chan BotReply)
		for _, event := range events {
			go b.dispatch(event, replyChannel)
		}
		for i := 0; i < len(events); i++ {
			go b.Provider.EventReplier(<-replyChannel)
		}
	}

	ctx.JSON(map[string]string{
		"message": "event dispatched",
	})
}

// New ...
func New(provider Provider) Bot {
	return Bot{
		Provider: provider,
		Commands: &Commands,
	}
}
