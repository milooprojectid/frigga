package bot

// Event ...
type Event struct {
	ID      string
	Message string
	Token   string
}

type eventReply struct {
	Messages []string
	Token    string
}

func (e Event) isTrigger() bool {
	return string(e.Message[0]) == "/"
}
