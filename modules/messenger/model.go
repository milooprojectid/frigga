package messenger

type sender struct {
	ID string `json:"id"`
}

type recipient struct {
	ID string `json:"id"`
}

// Message ...
type Message struct {
	Mid        string `json:"mid"`
	Text       string `json:"text"`
	QuickReply struct {
		Payload string `json:"payload"`
	} `json:"quick_reply"`
}

// TextMessage ...
type TextMessage struct {
	Sender    sender    `json:"sender"`
	Recipient recipient `json:"recipient"`
	Timestamp int64     `json:"timestamp"`
	Message   Message   `json:"message"`
}
