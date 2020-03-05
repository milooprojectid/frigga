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

// Messaging ...
type Messaging struct {
	Sender    sender    `json:"sender"`
	Recipient recipient `json:"recipient"`
	Timestamp int64     `json:"timestamp"`
	Message   Message   `json:"message"`
}

type entry struct {
	ID        string      `json:"id"`
	Time      int64       `json:"time"`
	Messaging []Messaging `json:"messaging"`
}

// Event ...
type Event struct {
	Object string  `json:"object"`
	Entry  []entry `json:"entry"`
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
