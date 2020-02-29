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
