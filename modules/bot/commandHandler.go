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

// Common
func startCommandTrigger(event BotEvent) ([]c.Message, error) {
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

	return c.GenerateTextMessages([]string{
		"Hi im MilooBot\n" + commandGreetMessage,
	}), nil
}

// Common
func helpCommandTrigger(event BotEvent) ([]c.Message, error) {
	commands, _ := repo.GetAllBotFeatures("active")
	templatePayload := map[string]interface{}{
		"Commands": commands,
	}
	message := template.ProcessFile("storage/help.tmpl", templatePayload)

	return c.GenerateTextMessages([]string{
		message,
	}), nil
}

// Common
func cancelCommandTrigger(event BotEvent) ([]c.Message, error) {
	ID := event.ID
	var message string

	if command, _ := repo.GetSession(ID); command == "" {
		message = commandInactiveMessage + " to cancel. I wasn't doing anything anyway. Zzzzz..."
	} else {
		message = commandCancelledMessage
		repo.UpdateSession(ID, "")
	}

	return c.GenerateTextMessages([]string{
		message,
		commandGreetMessage,
	}), nil
}

// Sentiment
func sentimentCommandTrigger(event BotEvent) ([]c.Message, error) {
	repo.UpdateSession(event.ID, "/sentiment")
	return c.GenerateTextMessages([]string{
		"Type the statement you want to analize",
	}), nil
}

// Sentiment
func sentimentCommandFeedback(event BotEvent, payload ...interface{}) ([]c.Message, error) {
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

	return c.GenerateTextMessages([]string{
		message,
	}), nil
}

// Summarization
func summarizeCommandTrigger(event BotEvent) ([]c.Message, error) {
	repo.UpdateSession(event.ID, "/summarize")
	return c.GenerateTextMessages([]string{
		"Type the statement or url you want to summarise",
	}), nil
}

// Summarization
func summarizeCommandFeedback(event BotEvent, payload ...interface{}) ([]c.Message, error) {
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

	return c.GenerateTextMessages([]string{
		message,
	}), nil
}

// Corona Report
func covid19CommandTrigger(event BotEvent) ([]c.Message, error) {
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

	return c.GenerateTextMessages(messages), nil
}

// Corona Subscriptions
func covid19SubscribeCommandTrigger(event BotEvent) ([]c.Message, error) {
	ID := event.ID
	provider := event.Provider

	messages := []string{
		"Thank you for subscribing, we will notify you about covid-19 daily",
	}
	cmd := "/subscov19"

	repo.SetCovid19SubsData(ID, provider)
	repo.LogSession(ID, cmd, "", "")

	return c.GenerateTextMessages(messages), nil
}

// Type Tester
func typeTesterCommandTrigger(event BotEvent) ([]c.Message, error) {
	return c.GenerateTextMessages([]string{
		"Ahoy sailor, youve come to test jaa ?",
	}), nil
}

// Type Tester
func typeTesterCommandFeedback(event BotEvent, payload ...interface{}) ([]c.Message, error) {
	var response c.Message
	var outType string = payload[0].(string)

	switch outType {
	case c.TextMessageType:
		{
			response = c.GenerateTextMessage("hi, i am a plain text, like your relationship :)")
		}

	case c.AudioMessageType:
		{
			response = c.GenerateAudioMessage("https://firebasestorage.googleapis.com/v0/b/miloo-phoenix.appspot.com/o/audio-example.mp3?alt=media&token=5684a144-6d09-42db-88d4-29c1445225ef")
		}

	case c.VideoMessageType:
		{
			response = c.GenerateVideoMessage("https://firebasestorage.googleapis.com/v0/b/miloo-phoenix.appspot.com/o/video-example.mp4?alt=media&token=13def274-3b0d-4771-9f15-6ab18064087b")
		}

	case c.ImageMessageType:
		{
			response = c.GenerateImageMessage("https://firebasestorage.googleapis.com/v0/b/miloo-phoenix.appspot.com/o/miloo-edited.png?alt=media&token=94ddd99a-d801-48c4-9345-98bffd91b264")
		}

	case c.LocationMessageType:
		{
			response = c.GenerateLocationMessage("-6.207644,106.809604")
		}

	}

	return []c.Message{
		response,
	}, nil
}
