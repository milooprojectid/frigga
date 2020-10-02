package bot

import (
	c "frigga/modules/common"
	repo "frigga/modules/repository"
	"frigga/modules/service"
)

// Commands containts all command available
var Commands commands

// Common Messaged
const (
	commandGreetMessage     = "You can see all available command with /help"
	commandFailedMessage    = "Hmm, sorry i have problem processing your command :("
	commandUnknownMessage   = "Hmm, i dont know that command ._."
	commandInactiveMessage  = "No active command"
	commandCancelledMessage = "Command cancelled"
)

// Command ...
type Command struct {
	Name        string
	Path        string
	Description string
	IsVisible   bool
	Trigger     func(event BotEvent) ([]c.Message, error)
	Feedback    func(event BotEvent, output ...interface{}) ([]c.Message, error)
}

func (c *Command) getEventType() string {
	var eventType string = "trigger"
	if c.Feedback == nil {
		eventType = "feedback"
	}
	return eventType
}

// Commands

type commands []Command

func (cs *commands) getCommand(path string) *Command {
	var command *Command

	for _, c := range Commands {
		if c.Path == path {
			command = &c
			break
		}
	}

	return command
}

func getRasaEventReply(token string, text string) BotReply {
	var response []struct {
		RecipientID string `json:"recipient_id"`
		Text        string `json:"text"`
	}
	service.All.CallSync("rasa", "predict", map[string]interface{}{"message": text}, &response)
	return BotReply{
		Messages: c.GenerateTextMessages([]string{response[0].Text}),
		Token:    token,
		Type:     "feedback",
	}
}

func getUnknownEventReply(token string) BotReply {
	return BotReply{
		Messages: c.GenerateTextMessages([]string{
			commandUnknownMessage,
			commandGreetMessage,
		}),
		Token: token,
		Type:  "feedback",
	}
}

// Execute ...
func (cs *commands) Execute(event BotEvent) BotReply {
	var command *Command
	var reply BotReply

	if event.isTrigger() {
		if command = cs.getCommand(event.Message.Text); command == nil {
			reply = getUnknownEventReply(event.Token)
		} else {
			eventType := command.getEventType()
			messages, _ := command.Trigger(event)
			reply = BotReply{
				Messages: messages,
				Token:    event.Token,
				Type:     eventType,
			}
		}
	} else if cmd, _ := repo.GetSession(event.ID); cmd == "" {
		reply = getRasaEventReply(event.Token, event.Message.Text)
	} else {
		command = cs.getCommand(cmd)
		messages, _ := command.Feedback(event, event.Message.Text)
		reply = BotReply{
			Messages: messages,
			Token:    event.Token,
			Type:     "feedback",
		}
	}

	return reply
}

// ExecuteInline ...
func (cs *commands) ExecuteInline(event BotEvent) BotReply {
	var reply BotReply
	commandPath, input, _ := event.getCommandAndInput()

	if command := cs.getCommand(commandPath); command == nil {
		reply = getUnknownEventReply(event.Token)
	} else if command.Feedback == nil {
		reply = BotReply{
			Messages: c.GenerateTextMessages([]string{
				commandFailedMessage,
			}),
			Token: event.Token,
			Type:  "feedback",
		}
	} else {
		messages, _ := command.Feedback(event, input)
		reply = BotReply{
			Messages: messages,
			Token:    event.Token,
			Type:     "feedback",
		}
	}

	return reply
}

// RegisterCommands ...
func RegisterCommands() {
	startCommand := Command{
		Name:     "Start",
		Path:     "/start",
		Trigger:  startCommandTrigger,
		Feedback: nil,
	}

	helpCommand := Command{
		Name:     "Help",
		Path:     "/help",
		Trigger:  helpCommandTrigger,
		Feedback: nil,
	}

	cancelCommand := Command{
		Name:     "Cancel",
		Path:     "/cancel",
		Trigger:  cancelCommandTrigger,
		Feedback: nil,
	}

	// sentimentCommand := Command{
	// 	Name:     "Sentiment",
	// 	Path:     "/sentiment",
	// 	Trigger:  sentimentCommandTrigger,
	// 	Feedback: sentimentCommandFeedback,
	// }

	summarizeCommand := Command{
		Name:     "Summarization",
		Path:     "/summarize",
		Trigger:  summarizeCommandTrigger,
		Feedback: summarizeCommandFeedback,
	}

	covid19ReportCommand := Command{
		Name:     "Covid19Report",
		Path:     "/corona",
		Trigger:  covid19CommandTrigger,
		Feedback: nil,
	}

	covid19SubscriptionCommand := Command{
		Name:     "Covid19Subs",
		Path:     "/subscov19",
		Trigger:  covid19SubscribeCommandTrigger,
		Feedback: nil,
	}

	aboutCommand := Command{
		Name:     "About",
		Path:     "/about",
		Trigger:  aboutCommandTrigger,
		Feedback: nil,
	}

	hadistCommand := Command{
		Name:     "HadistRetrieval",
		Path:     "/hadist",
		Trigger:  hadistCommandTrigger,
		Feedback: hadistCommandFeedback,
	}

	// initialize to singletons
	Commands = commands{
		startCommand,
		helpCommand,
		cancelCommand,
		// sentimentCommand,
		summarizeCommand,
		covid19ReportCommand,
		covid19SubscriptionCommand,
		aboutCommand,
		hadistCommand,
	}
}
