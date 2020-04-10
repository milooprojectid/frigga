package line

// Name ...
const Name = "line"

// Event ...
type Event struct {
	ReplyToken string `json:"replyToken"`
	Type       string `json:"type"`
	Mode       string `json:"mode"`
	Timestamp  int64  `json:"timestamp"`
	Source     struct {
		Type   string `json:"type"`
		UserID string `json:"userId"`
	} `json:"source"`
	Message struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"message"`
}

// ItemAction ...
type ItemAction struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	Text  string `json:"text"`
}

// QuickReplyItem ...
type QuickReplyItem struct {
	Type     string     `json:"type"`
	ImageURL *string    `json:"imageUrl,omitempty"`
	Action   ItemAction `json:"action"`
}

// QuickReply ...
type QuickReply struct {
	Items []QuickReplyItem `json:"items"`
}

// TextReplyMessage ...
type TextReplyMessage struct {
	Type       string      `json:"type"`
	Text       string      `json:"text"`
	QuickReply *QuickReply `json:"quickReply,omitempty"`
}

// AudioReplyMessage ...
type AudioReplyMessage struct {
	Type               string      `json:"type"`
	OriginalContentURL string      `json:"originalContentUrl"`
	Duration           int         `json:"duration,omitempty"`
	QuickReply         *QuickReply `json:"quickReply,omitempty"`
}

// VideoReplyMessage ...
type VideoReplyMessage struct {
	Type               string      `json:"type"`
	OriginalContentURL string      `json:"originalContentUrl"`
	PreviewImageURL    string      `json:"previewImageUrl,omitempty"`
	QuickReply         *QuickReply `json:"quickReply,omitempty"`
}

// ImageReplyMessage ...
type ImageReplyMessage struct {
	Type               string      `json:"type"`
	OriginalContentURL string      `json:"originalContentUrl"`
	PreviewImageURL    string      `json:"previewImageUrl,omitempty"`
	QuickReply         *QuickReply `json:"quickReply,omitempty"`
}

// LocationReplyMessage ...
type LocationReplyMessage struct {
	Type       string      `json:"type"`
	Title      string      `json:"title"`
	Address    string      `json:"address"`
	Latitude   float64     `json:"latitude"`
	Longitude  float64     `json:"longitude"`
	QuickReply *QuickReply `json:"quickReply,omitempty"`
}

// Request ...
type Request struct {
	Destination string  `json:"destination"`
	Events      []Event `json:"events"`
}

type userProfile struct {
	DisplayName   string `json:"displayName"`
	UserID        string `json:"userId"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}
