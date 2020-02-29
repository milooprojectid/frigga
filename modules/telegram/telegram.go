package telegram

import (
	"bytes"
	"encoding/json"
	"net/http"

	// repo "frigga/modules/repository"
	// service "frigga/modules/service"

	"github.com/kataras/iris"
)

func (e Event) isCommand() bool {
	return string(e.Message.Text[0]) == "/"
}

// Process ...
func (e *Event) Process(c chan messageReply) {
	// var message string
	// var input string = e.Message.Text
	// var chatID string = strconv.Itoa(e.Message.Chat.ID)

	// switch input {
	// case "/start":
	// 	{
	// 		message = "Hi im Miloo\n" + commands
	// 		repo.InitSession(chatID, e.Message.Chat.FirstName, e.Message.Chat.LastName)
	// 	}
	// case "/sentiment":
	// 	{
	// 		message = "Type the statement you want to analize"
	// 		repo.UpdateSession(chatID, input)
	// 	}
	// case "/summarize":
	// 	{
	// 		message = "Type the statement or url you want to summarise"
	// 		repo.UpdateSession(chatID, input)
	// 	}
	// case "/cancel":
	// 	{
	// 		if command, _ := repo.GetSession(chatID); command == "" {
	// 			message = "No active command to cancel. I wasn't doing anything anyway. Zzzzz..."
	// 		} else {
	// 			message = "Command cancelled"
	// 			repo.UpdateSession(chatID, "")
	// 		}
	// 	}
	// default:
	// 	{
	// 		message = ""
	// 	}
	// }

	// if message == "" {
	// 	if len(input) == 0 {
	// 		message = "You have to type something ._."
	// 	} else if string(input[0]) == "/" {
	// 		message = "I cant understand that command\n" + commands
	// 	}

	// 	if command, _ := repo.GetSession(chatID); command == "" {
	// 		message = "Hmmm, theres no active command >_<"
	// 	} else if command == "/summarize" {
	// 		var result summarizationResult
	// 		service.All.CallSync("storm", "summarizeText", map[string]string{"text": input}, &result)
	// 		message = result.Data.Summary
	// 		repo.LogHistory(chatID, command, input, message)
	// 	} else if command == "/sentiment" {
	// 		var result sentimentResult
	// 		service.All.CallSync("morbius", "sentiment", map[string]string{"text": input}, &result)
	// 		message = result.Data.Description
	// 		repo.LogHistory(chatID, command, input, message)
	// 	}

	// 	repo.UpdateSession(chatID, "")
	// }

	// c <- messageReply{
	// 	ChatID:  e.Message.Chat.ID,
	// 	Message: message,
	// }
}

// Handler ...
func (b *Bot) Handler(ctx iris.Context) {
	var event Event
	if err := ctx.ReadJSON(&event); err != nil {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.WriteString(err.Error())
		return
	}

	replies := make(chan messageReply)
	go event.Process(replies)

	r := <-replies
	go b.Reply(r.ChatID, r.Message)

	ctx.JSON(event)
	return
}

// Reply ...
func (b *Bot) Reply(chatID int, message string) error {
	payload := map[string]interface{}{"chat_id": chatID, "text": message}
	requestBody, _ := json.Marshal(payload)
	if _, err := http.Post(apiURL+b.Token+"/sendMessage", "application/json", bytes.NewBuffer(requestBody)); err != nil {
		return err
	}
	return nil
}

// New returns new bot instance
func New(token string) Bot {
	return Bot{
		Token: token,
	}
}
