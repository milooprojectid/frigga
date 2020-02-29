package bot

// Event ...
type Event struct {
	ChatID     string
	Message    string
	ReplyToken string
}

func (e Event) isCommand() bool {
	return string(e.Message[0]) == "/"
}
