package bot

import (
	"context"
	line "frigga/modules/line"
	messenger "frigga/modules/messenger"
	repo "frigga/modules/repository"
	service "frigga/modules/service"
	storm "frigga/modules/service/storm"
	telegram "frigga/modules/telegram"
	"frigga/modules/template"
	"time"
)

// Commands containts all command available
var Commands commands

const commandGreetMessage = "You can control me by sending these commands"
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

func getUnknownEventReply(token string) eventReply {
	return eventReply{
		Messages: []string{"I dont know that command ._.", commandGreetMessage},
		Token:    token,
		Type:     "feedback",
	}
}

func (cs *commands) execute(event Event, provider string) eventReply {
	var command *Command
	var reply eventReply

	if event.isTrigger() {
		if command = cs.getCommand(event.Message); command == nil {
			reply = getUnknownEventReply(event.Token)
		} else {
			var eventType string = "trigger"
			if command.Feedback == nil {
				eventType = "feedback"
			}
			messages, _ := command.Trigger(event.ID, provider)
			reply = eventReply{
				Messages: messages,
				Token:    event.Token,
				Type:     eventType,
			}
		}
	} else if cmd, _ := repo.GetSession(event.ID); cmd == "" {
		reply = eventReply{
			Messages: []string{"No active command", commandGreetMessage},
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

func (cs *commands) executeInline(event Event, provider string) eventReply {
	var reply eventReply
	commandPath, input, _ := event.getCommandAndInput()

	if command := cs.getCommand(commandPath); command == nil {
		reply = getUnknownEventReply(event.Token)
	} else {
		messages, _ := command.Feedback(event.ID, input)
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
		Trigger: covid19Command,
	}

	covid19SubscriptionCommand := Command{
		Name:    "Covid19Subs",
		Path:    "/subscov19",
		Trigger: covid19SubscribeCommand,
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

func startCommandTrigger(payload ...interface{}) ([]string, error) {
	var ID string = payload[0].(string)
	var provider = payload[1].(string)
	var getter func(id string) (string, error)

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

func helpCommandTrigger(payload ...interface{}) ([]string, error) {
	return []string{commandGreetMessage}, nil
}

func cancelCommandTrigger(payload ...interface{}) ([]string, error) {
	var ID string = payload[0].(string)
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

func sentimentCommandTrigger(payload ...interface{}) ([]string, error) {
	var ID string = payload[0].(string)

	repo.UpdateSession(ID, "/sentiment")
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

	repo.LogSession(ID, cmd, input, message)
	repo.UpdateSession(ID, "")

	return []string{
		message,
	}, nil
}

func summarizeCommandTrigger(payload ...interface{}) ([]string, error) {
	var ID string = payload[0].(string)

	repo.UpdateSession(ID, "/summarize")
	return []string{
		"Type the statement or url you want to summarise",
	}, nil
}

func summarizeCommandFeedback(payload ...interface{}) ([]string, error) {
	var ID string = payload[0].(string)
	var input string = payload[1].(string)

	var message string
	var result string
	cmd := "/summarize"

	request := &storm.SummarizeRequest{
		Text: input,
	}
	response, _ := (*storm.Client).Summarize(context.Background(), request)
	result = response.GetSummary()

	if result == "" {
		message = commandFailedMessage
	} else {
		message = result
	}

	repo.LogSession(ID, cmd, input, message)
	repo.UpdateSession(ID, "")

	return []string{
		message,
	}, nil
}

func covid19Command(payload ...interface{}) ([]string, error) {
	var ID string = payload[0].(string)
	cmd := "/corona"

	data, _ := repo.GetCovid19Data()
	templatePayload := map[string]interface{}{
		"Date":      time.Now().Format("02 Mar 2006"),
		"Confirmed": data.Confirmed,
		"Recovered": data.Recovered,
		"Deceased":  data.Deceased,
	}

	message := template.ProcessFile("storage/covid19.tmpl", templatePayload)
	messages := []string{
		message,
	}

	repo.LogSession(ID, cmd, "", "")

	return messages, nil
}

func covid19SubscribeCommand(payload ...interface{}) ([]string, error) {
	var ID string = payload[0].(string)
	var provider string = payload[1].(string)

	messages := []string{
		"Thank you for subscribing, we will notify you about covid-19 daily",
	}
	cmd := "/subscov19"

	repo.SetCovid19SubsData(ID, provider)
	repo.LogSession(ID, cmd, "", "")

	return messages, nil
}
