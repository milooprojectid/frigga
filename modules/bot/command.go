package bot

import (
	"context"
	c "frigga/modules/common"
	line "frigga/modules/providers/line"
	messenger "frigga/modules/providers/messenger"
	telegram "frigga/modules/providers/telegram"
	repo "frigga/modules/repository"
	morbius "frigga/modules/service/morbius"
	storm "frigga/modules/service/storm"
	"frigga/modules/template"
	"time"
)

// Commands containts all command available
var Commands commands

const commandGreetMessage = "You can control me by sending these commands"
const commandFailedMessage = "Hmm, sorry i have problem processing your message :("

// Command ...
type Command struct {
	Name     string
	Path     string
	Trigger  func(event Event) ([]string, error)
	Feedback func(event Event, output ...interface{}) ([]string, error)
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

func getUnknownEventReply(token string) eventReply {
	return eventReply{
		Messages: c.GenerateTextMessages([]string{
			"I dont know that command ._.",
			commandGreetMessage,
		}),
		Token: token,
		Type:  "feedback",
	}
}

func (cs *commands) execute(event Event, provider string) eventReply {
	var command *Command
	var reply eventReply

	if event.isTrigger() {
		if command = cs.getCommand(event.Message.Text); command == nil {
			reply = getUnknownEventReply(event.Token)
		} else {
			var eventType string = "trigger"
			if command.Feedback == nil {
				eventType = "feedback"
			}
			messages, _ := command.Trigger(event)
			reply = eventReply{
				Messages: c.GenerateTextMessages(messages),
				Token:    event.Token,
				Type:     eventType,
			}
		}
	} else if cmd, _ := repo.GetSession(event.ID); cmd == "" {
		reply = eventReply{
			Messages: c.GenerateTextMessages([]string{
				"No active command",
				commandGreetMessage,
			}),
			Token: event.Token,
			Type:  "feedback",
		}
	} else {
		command = cs.getCommand(cmd)
		messages, _ := command.Feedback(event, event.Message.Text)
		reply = eventReply{
			Messages: c.GenerateTextMessages(messages),
			Token:    event.Token,
			Type:     "feedback",
		}
	}

	return reply
}

func (cs *commands) executeInline(event Event, provider string) eventReply {
	var reply eventReply
	commandPath, input, _ := event.getCommandAndInput()

	if command := cs.getCommand(commandPath); command == nil {
		reply = getUnknownEventReply(event.Token)
	} else {
		messages, _ := command.Feedback(event, input)
		reply = eventReply{
			Messages: c.GenerateTextMessages(messages),
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

	helpCommand := Command{
		Name:    "Help",
		Path:    "/help",
		Trigger: helpCommandTrigger,
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

	covid19SummaryCommand := Command{
		Name:    "Covid19",
		Path:    "/corona",
		Trigger: covid19CommandTrigger,
	}

	covid19SubscriptionCommand := Command{
		Name:    "Covid19Subs",
		Path:    "/subscov19",
		Trigger: covid19SubscribeCommandTrigger,
	}

	// initialize to singletons
	Commands = commands{
		startCommand,
		helpCommand,
		cancelCommand,
		sentimentCommand,
		summarizeCommand,
		covid19SummaryCommand,
		covid19SubscriptionCommand,
	}
}

// explicit handler

func startCommandTrigger(event Event) ([]string, error) {
	var getter func(id string) (string, error)

	ID := event.ID
	provider := event.Provider

	switch provider {
	case "telegram":
		getter = telegram.GetUserName
	case "line":
		getter = line.GetUserName
	case "messenger":
		getter = messenger.GetUserName
	default:
		getter = func(id string) (string, error) {
			return id, nil
		}
	}

	name, _ := getter(ID)
	repo.InitSession(provider, ID, name)

	return []string{
		"Hi im Miloo\n" + commandGreetMessage,
	}, nil
}

func helpCommandTrigger(event Event) ([]string, error) {
	return []string{commandGreetMessage}, nil
}

func cancelCommandTrigger(event Event) ([]string, error) {
	ID := event.ID
	var message string

	if command, _ := repo.GetSession(ID); command == "" {
		message = "No active command to cancel. I wasn't doing anything anyway. Zzzzz..."
	} else {
		message = "Command cancelled"
		repo.UpdateSession(ID, "")
	}

	return []string{
		message,
		commandGreetMessage,
	}, nil
}

func sentimentCommandTrigger(event Event) ([]string, error) {
	repo.UpdateSession(event.ID, "/sentiment")
	return []string{
		"Type the statement you want to analize",
	}, nil
}

func sentimentCommandFeedback(event Event, payload ...interface{}) ([]string, error) {
	ID := event.ID
	input := payload[0].(string)

	var message string
	cmd := "/sentiment"

	request := &morbius.SentimentRequest{
		Text: input,
	}

	if response, err := (*morbius.Client).Sentiment(context.Background(), request); err != nil || response.GetDescription() == "" {
		message = commandFailedMessage
	} else {
		message = response.GetDescription()
	}

	repo.LogSession(ID, cmd, input, message)
	repo.UpdateSession(ID, "")

	return []string{
		message,
	}, nil
}

func summarizeCommandTrigger(event Event) ([]string, error) {
	repo.UpdateSession(event.ID, "/summarize")
	return []string{
		"Type the statement or url you want to summarise",
	}, nil
}

func summarizeCommandFeedback(event Event, payload ...interface{}) ([]string, error) {
	ID := event.ID
	input := payload[0].(string)

	var message string
	cmd := "/summarize"

	request := &storm.SummarizeRequest{
		Text: input,
	}

	if response, err := (*storm.Client).Summarize(context.Background(), request); err != nil || response.GetSummary() == "" {
		message = commandFailedMessage
	} else {
		message = response.GetSummary()
	}

	repo.LogSession(ID, cmd, input, message)
	repo.UpdateSession(ID, "")

	return []string{
		message,
	}, nil
}

func covid19CommandTrigger(event Event) ([]string, error) {
	cmd := "/corona"

	data, _ := repo.GetCovid19Data()
	templatePayload := map[string]interface{}{
		"Date":      time.Now().Format("02 Jan 2006"),
		"Confirmed": data.Confirmed,
		"Recovered": data.Recovered,
		"Deceased":  data.Deceased,
		"Treated":   data.Confirmed - data.Recovered - data.Deceased,
	}

	message := template.ProcessFile("storage/covid19.tmpl", templatePayload)
	messages := []string{
		message,
	}

	repo.LogSession(event.ID, cmd, "", "")

	return messages, nil
}

func covid19SubscribeCommandTrigger(event Event) ([]string, error) {
	ID := event.ID
	provider := event.Provider

	messages := []string{
		"Thank you for subscribing, we will notify you about covid-19 daily",
	}
	cmd := "/subscov19"

	repo.SetCovid19SubsData(ID, provider)
	repo.LogSession(ID, cmd, "", "")

	return messages, nil
}
