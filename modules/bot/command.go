package bot

import (
	"frigga/modules/service"
)

// Commands containts all command available
var Commands commands

const commandText = "You can control me by sending these commands:	\n/sentiment - run sentiment analysis on a text\n/summarize - summarise a content of a link or text\n/cancel - terminate currently running command"

type commandhandler func() ([]string, error)

// Command ...
type Command struct {
	Name     string
	Path     string
	Trigger  func(param ...interface{}) ([]string, error)
	Feedback func(param ...interface{}) ([]string, error)
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

func (cs *commands) execute(event Event) eventReply {
	var command *Command
	var reply eventReply

	if event.isTrigger() {
		if command = cs.getCommand(event.Message); command == nil {
			reply = eventReply{"I dont know that command ._.", event.Token}
		} else {
			messages, _ := command.Trigger(event.ID)
			reply = eventReply{messages[0], event.Token}
		}

	} else if cmd, _ := GetSession(event.ID); cmd == "" {
		reply = eventReply{"No active command", event.Token}
	} else {
		command = cs.getCommand(cmd)
		messages, _ := command.Feedback(event.ID, event.Message)
		reply = eventReply{messages[0], event.Token}
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

	// initialize to singletons
	Commands = commands{
		startCommand,
		cancelCommand,
		sentimentCommand,
	}
}

// explicit handler

func startCommandTrigger(param ...interface{}) ([]string, error) {
	ID := param[0].(string)
	InitSession(ID)
	return []string{
		"Hi im Miloo\n" + commandText,
	}, nil
}

func cancelCommandTrigger(param ...interface{}) ([]string, error) {
	var message string
	ID := param[0].(string)

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

func sentimentCommandTrigger(param ...interface{}) ([]string, error) {
	ID := param[0].(string)
	UpdateSession(ID, "/sentiment")
	return []string{
		"Type the statement you want to analize",
	}, nil
}

func sentimentCommandFeedback(param ...interface{}) ([]string, error) {
	var result service.SentimentResult
	ID := param[0].(string)
	input := param[1].(string)
	cmd := "/sentiment"

	service.All.CallSync("morbius", "sentiment", map[string]string{"text": input}, &result)
	message := result.Data.Description

	LogSession(ID, cmd, input, message)
	UpdateSession(ID, "")

	return []string{
		message,
	}, nil
}
