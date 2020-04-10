package messenger

// Name ...
const Name = "messenger"

type sender struct {
	ID string `json:"id"`
}

type recipient struct {
	ID string `json:"id"`
}

// QuickReply ...
type QuickReply struct {
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	ImageURL    string `json:"image_url,omitempty"`
	Payload     string `json:"payload"`
}

// Message ...
type Message struct {
	Mid  string `json:"mid"`
	Text string `json:"text"`
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

// BotEvent ...
type BotEvent struct {
	Object string  `json:"object"`
	Entry  []entry `json:"entry"`
}

// SendPayloadRecipient ...
type SendPayloadRecipient struct {
	ID string `json:"id"`
}

// SendPayloadMessage ...
type SendPayloadMessage struct {
	Text       string        `json:"text"`
	QuickReply *[]QuickReply `json:"quick_replies,omitempty"`
}

// SendPayload ...
type SendPayload struct {
	MessagingType string               `json:"messaging_type"`
	Recipient     SendPayloadRecipient `json:"recipient"`
	Message       SendPayloadMessage   `json:"message"`
}

type userProfile struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ProfilePic string `json:"profile_pic"`
	Locale     string `json:"locale"`
	Timezone   int    `json:"timezone"`
	Gender     string `json:"gender"`
}
