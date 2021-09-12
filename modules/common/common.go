package common

// TextMessageType ..
// AudioMessageType ..
// VideoMessageType ..
// ImageMessageType ..
// LocationMessageType ..
// AlbumMessageType ..
const (
	TextMessageType     string = "text"
	AudioMessageType    string = "audio"
	VideoMessageType    string = "video"
	ImageMessageType    string = "image"
	LocationMessageType string = "location"
	AlbumMessageType    string = "album"
)

type AlbumItem struct {
	Type string `json:"type" validate:"required,min=4"`
	Body string `json:"body" validate:"required,min=3"`
}

// Message ..
type Message struct {
	Type  string       `json:"type" validate:"required"`
	Text  string       `json:"text" validate:"required"`
	Album *[]AlbumItem `json:"album,omitempty"`
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

// GenerateAlbumMessage
func GenerateAlbumMessage(text string, albumItems []AlbumItem) Message {
	return Message{
		Text:  text,
		Type:  AlbumMessageType,
		Album: &albumItems,
	}
}
