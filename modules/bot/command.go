package bot

const commandText = "You can control me by sending these commands:	\n/sentiment - run sentiment analysis on a text\n/summarize - summarise a content of a link or text\n/cancel - terminate currently running command"

// Commands containts all command available
var Commands commands

type commandhandler func() ([]string, error)

// Command ...
type Command struct {
	Name     string
	Path     string
	Trigger  func(param ...interface{}) ([]string, error)
	Feedback func(param ...interface{}) ([]string, error)
}

type commands []Command

func (cs *commands) execute(event Event) eventReply {
	var command *Command
	var reply eventReply

	for _, c := range Commands {
		if c.Path == event.Message {
			command = &c
		}
	}

	if command == nil {
		reply = eventReply{"I cant understand that command ._.", event.Token}
	}

	if event.isTrigger() {
		messages, _ := command.Trigger(event.ID)
		reply = eventReply{messages[0], event.Token}
	} else {
		messages, _ := command.Feedback(event.ID)
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

	// initialize to singletons
	Commands = commands{
		startCommand,
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
