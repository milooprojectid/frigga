package line

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
	Type   string     `json:"type"`
	Action ItemAction `json:"action"`
}

// ImageURL *string    `json:"imageUrl,omitempty"`

// QuickReply ...
type QuickReply struct {
	Items []QuickReplyItem `json:"items"`
}

// ReplyMessage ...
type ReplyMessage struct {
	Type       string     `json:"type"`
	Text       string     `json:"text"`
	QuickReply QuickReply `json:"quickReply,omitempty"`
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
