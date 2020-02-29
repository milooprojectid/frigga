package bot

const commandText = "You can control me by sending these commands:\n/sentiment - run sentiment analysis on a text\n/summarize - summarise a content of a link or text\n/cancel - terminate currently running command"

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

type commands map[string]Command

// RegisterCommands ...
func RegisterCommands() {
	startCommand := Command{
		Name:    "Start",
		Path:    "/start",
		Trigger: startCommandTrigger,
	}

	// initiaze to singletons
	Commands = commands{
		"start": startCommand,
	}
}

// explicit handler

func startCommandTrigger(param ...interface{}) ([]string, error) {
	return []string{
		"Hi im Miloo\n" + commandText,
	}, nil
}
