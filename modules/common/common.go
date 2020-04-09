package common

// TextMessageType ..
// AudioMessageType ..
// VideoMessageType ..
// ImageMessageType ..
// LocationMessageType ..
const (
	TextMessageType     = "text"
	AudioMessageType    = "audio"
	VideoMessageType    = "video"
	ImageMessageType    = "image"
	LocationMessageType = "location"
)

// Message ..
type Message struct {
	Type string
	Text string
}

// GenerateTextMessages ...
func GenerateTextMessages(texts []string) []Message {
	var messages []Message
	for _, text := range texts {
		messages = append(messages, GenerateTextMessage(text))
	}
	return messages
}

// GenerateTextMessage ...
func GenerateTextMessage(text string) Message {
	return Message{
		Text: text,
		Type: TextMessageType,
	}
}
