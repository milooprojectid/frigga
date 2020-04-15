package common

// TextMessageType ..
// AudioMessageType ..
// VideoMessageType ..
// ImageMessageType ..
// LocationMessageType ..
const (
	TextMessageType     string = "text"
	AudioMessageType    string = "audio"
	VideoMessageType    string = "video"
	ImageMessageType    string = "image"
	LocationMessageType string = "location"
)

// Message ..
type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
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

// GenerateVideoMessage ...
func GenerateVideoMessage(url string) Message {
	return Message{
		Text: url,
		Type: VideoMessageType,
	}
}

// GenerateAudioMessage ...
func GenerateAudioMessage(url string) Message {
	return Message{
		Text: url,
		Type: AudioMessageType,
	}
}

// GenerateImageMessage ...
func GenerateImageMessage(url string) Message {
	return Message{
		Text: url,
		Type: ImageMessageType,
	}
}

// GenerateLocationMessage ...
func GenerateLocationMessage(location string) Message {
	return Message{
		Text: location,
		Type: LocationMessageType,
	}
}
