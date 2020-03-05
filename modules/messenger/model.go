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

// Entry ...
type Entry struct {
	Sender    sender    `json:"sender"`
	Recipient recipient `json:"recipient"`
	Timestamp int64     `json:"timestamp"`
	Message   Message   `json:"message"`
}

// Event ...
type Event struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

// SendPayloadRecipient ...
type SendPayloadRecipient struct {
	ID string `json:"id"`
}

// SendPayloadMessage ...
type SendPayloadMessage struct {
	Text string `json:"text"`
}

// SendPayload ...
type SendPayload struct {
	MessagingType string               `json:"messaging_type"`
	Recipient     SendPayloadRecipient `json:"recipient"`
	Message       SendPayloadMessage   `json:"message"`
}
