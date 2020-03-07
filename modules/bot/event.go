package bot

import (
	"errors"
	"strings"
)

// Event ...
type Event struct {
	ID      string
	Message string
	Token   string
}

type eventReply struct {
	Messages []string
	Token    string
	Type     string
}

func (e Event) isTrigger() bool {
	return e.Message != "" && string(e.Message[0]) == "/"
}

func (e Event) isInline() bool {
	var message string = e.Message

	if message == "" {
		return false
	}

	if string(message[0]) != "/" {
		return false
	}

	splitted := strings.SplitN(message, " ", 2)
	if len(splitted) != 2 {
		return false
	}

	return true
}

func (e Event) getCommandAndInput() (string, string, error) {
	var message string = e.Message
	var command string
	var input string

	if message == "" {
		return command, input, errors.New("message is empty")
	}

	if string(message[0]) != "/" {
		return command, input, errors.New("not a command")
	}

	splitted := strings.SplitN(message, " ", 2)
	if len(splitted) != 2 {
		return command, input, errors.New("cant split")
	}

	command = splitted[0]
	input = splitted[1]

	return command, input, nil
}

func (e eventReply) isTrigger() bool {
	return e.Type == "trigger"
}
