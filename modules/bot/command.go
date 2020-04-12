package bot

import (
	c "frigga/modules/common"
	repo "frigga/modules/repository"
)

// Commands containts all command available
var Commands commands

// Common Messaged
const (
	commandGreetMessage     = "You can control me by sending these commands below"
	commandFailedMessage    = "Hmm, sorry i have problem processing your command :("
	commandUnknownMessage   = "I dont know that command ._."
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

func (cs *commands) execute(event BotEvent, provider string) BotReply {
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
		reply = BotReply{
			Messages: c.GenerateTextMessages([]string{
				commandInactiveMessage,
				commandGreetMessage,
			}),
			Token: event.Token,
			Type:  "feedback",
		}
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

func (cs *commands) executeInline(event BotEvent, provider string) BotReply {
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
		Name:        "Start",
		Path:        "/start",
		Description: "initiate bot session",
		IsVisible:   false,
		Trigger:     startCommandTrigger,
		Feedback:    nil,
	}

	helpCommand := Command{
		Name:        "Help",
		Path:        "/help",
		Description: "return all available commands",
		IsVisible:   true,
		Trigger:     helpCommandTrigger,
		Feedback:    nil,
	}

	cancelCommand := Command{
		Name:        "Cancel",
		Path:        "/cancel",
		Description: "cancel ongoing command session",
		IsVisible:   false,
		Trigger:     cancelCommandTrigger,
		Feedback:    nil,
	}

	sentimentCommand := Command{
		Name:        "Sentiment",
		Path:        "/sentiment",
		Description: "analize sentiment of given text",
		IsVisible:   true,
		Trigger:     sentimentCommandTrigger,
		Feedback:    sentimentCommandFeedback,
	}

	summarizeCommand := Command{
		Name:        "Summarization",
		Path:        "/summarize",
		Description: "summarize content of text or news url",
		IsVisible:   true,
		Trigger:     summarizeCommandTrigger,
		Feedback:    summarizeCommandFeedback,
	}

	covid19ReportCommand := Command{
		Name:        "Covid19Report",
		Path:        "/corona",
		Description: "show current covid-19 progress in indonesia",
		IsVisible:   true,
		Trigger:     covid19CommandTrigger,
		Feedback:    nil,
	}

	covid19SubscriptionCommand := Command{
		Name:        "Covid19Subs",
		Path:        "/subscov19",
		Description: "subscribe to covid-19 progress in indonesia daily notification",
		IsVisible:   true,
		Trigger:     covid19SubscribeCommandTrigger,
		Feedback:    nil,
	}

	// initialize to singletons
	Commands = commands{
		startCommand,
		helpCommand,
		cancelCommand,
		sentimentCommand,
		summarizeCommand,
		covid19ReportCommand,
		covid19SubscriptionCommand,
	}
}

// GetAvailableCommandsPath ...
func GetAvailableCommandsPath() []string {
	var paths []string

	for _, command := range Commands {
		if command.IsVisible {
			paths = append(paths, command.Path)
		}
	}

	return paths
}
