package bot

import (
	line "frigga/modules/line"
	"frigga/modules/service"
	telegram "frigga/modules/telegram"
)

// Commands containts all command available
var Commands commands

const commandGreetMessage = "You can control me by sending these commands:	\n/sentiment - run sentiment analysis on a text\n/summarize - summarise a content of a link or text\n/cancel - terminate currently running command"
const commandFailedMessage = "Hmm, sorry i have problem processing your message :("

type commandhandler func(payload ...interface{}) ([]string, error)

// Command ...
type Command struct {
	Name     string
	Path     string
	Trigger  commandhandler
	Feedback commandhandler
}

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

func (cs *commands) execute(event Event, provider string) eventReply {
	var command *Command
	var reply eventReply

	if event.isTrigger() {
		if command = cs.getCommand(event.Message); command == nil {
			reply = eventReply{
				Messages: []string{"I dont know that command ._."},
				Token:    event.Token,
				Type:     "feedback",
			}
		} else {
			messages, _ := command.Trigger(event.ID, provider)
			reply = eventReply{
				Messages: messages,
				Token:    event.Token,
				Type:     "trigger",
			}
		}

	} else if cmd, _ := GetSession(event.ID); cmd == "" {
		reply = eventReply{
			Messages: []string{"No active command"},
			Token:    event.Token,
			Type:     "feedback",
		}
	} else {
		command = cs.getCommand(cmd)
		messages, _ := command.Feedback(event.ID, event.Message)
		reply = eventReply{
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
		Name:    "Start",
		Path:    "/start",
		Trigger: startCommandTrigger,
	}

	cancelCommand := Command{
		Name:    "Cancel",
		Path:    "/cancel",
		Trigger: cancelCommandTrigger,
	}

	sentimentCommand := Command{
		Name:     "Sentiment",
		Path:     "/sentiment",
		Trigger:  sentimentCommandTrigger,
		Feedback: sentimentCommandFeedback,
	}

	summarizeCommand := Command{
		Name:     "Summarization",
		Path:     "/summarize",
		Trigger:  summarizeCommandTrigger,
		Feedback: summarizeCommandFeedback,
	}

	// initialize to singletons
	Commands = commands{
		startCommand,
		cancelCommand,
		sentimentCommand,
		summarizeCommand,
	}
}

// explicit handler

func startCommandTrigger(payload ...interface{}) ([]string, error) {
	var ID string = payload[0].(string)
	var provider = payload[1].(string)
	var getter func(id string) (string, error)

	switch provider {
	case "telegram":
		getter = telegram.GetUserName
	case "line":
		getter = line.GetUserName
	}

	name, _ := getter(ID)
	InitSession(ID, name)

	return []string{
		"Hi im Miloo\n" + commandGreetMessage,
	}, nil
}

func cancelCommandTrigger(payload ...interface{}) ([]string, error) {
	var ID string = payload[0].(string)
	var message string

	if command, _ := GetSession(ID); command == "" {
		message = "No active command to cancel. I wasn't doing anything anyway. Zzzzz..."
	} else {
		message = "Command cancelled"
		UpdateSession(ID, "")
	}

	return []string{
		message,
	}, nil
}

func sentimentCommandTrigger(payload ...interface{}) ([]string, error) {
	var ID string = payload[0].(string)

	UpdateSession(ID, "/sentiment")
	return []string{
		"Type the statement you want to analize",
	}, nil
}

func sentimentCommandFeedback(payload ...interface{}) ([]string, error) {
	var ID string = payload[0].(string)
	var input string = payload[1].(string)

	var message string
	var result service.SentimentResult
	cmd := "/sentiment"

	if err := service.All.CallSync("morbius", "sentiment", map[string]string{"text": input}, &result); err != nil {
		message = commandFailedMessage
	} else if result.Data.Description == "" {
		message = commandFailedMessage
	} else {
		message = result.Data.Description
	}

	LogSession(ID, cmd, input, message)
	UpdateSession(ID, "")

	return []string{
		message,
	}, nil
}

func summarizeCommandTrigger(payload ...interface{}) ([]string, error) {
	var ID string = payload[0].(string)

	UpdateSession(ID, "/summarize")
	return []string{
		"Type the statement or url you want to summarise",
	}, nil
}

func summarizeCommandFeedback(payload ...interface{}) ([]string, error) {
	var ID string = payload[0].(string)
	var input string = payload[1].(string)

	var message string
	var result service.SummarizationResult
	cmd := "/summarize"

	if err := service.All.CallSync("storm", "summarizeText", map[string]string{"text": input}, &result); err != nil {
		message = commandFailedMessage
	} else if result.Data.Summary == "" {
		message = commandFailedMessage
	} else {
		message = result.Data.Summary
	}

	LogSession(ID, cmd, input, message)
	UpdateSession(ID, "")

	return []string{
		message,
	}, nil
}
