package telegram

const apiURL = "https://api.telegram.org/bot"
const commands = "You can control me by sending these commands:\n/sentiment - run sentiment analysis on a text\n/summarize - summarise a content of a link or text\n/cancel - terminate currently running command"

type sentimentResult struct {
	Data struct {
		Class       int    `json:"class"`
		Description string `json:"description"`
	} `json:"data"`
	Message string `json:"message"`
}

type summarizationResult struct {
	Data struct {
		Summary string `json:"summary"`
	} `json:"data"`
	Message string `json:"message"`
}

// Bot ...
type Bot struct {
	Token string
}

type user struct {
	ID        int    `json:"id"`
	IsBOT     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

type chat struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type message struct {
	MessageID int    `json:"message_id"`
	From      user   `json:"from"`
	Date      int    `json:"date"`
	Chat      chat   `json:"chat"`
	Text      string `json:"text,omitempty"`
}

// Update ...
type Update struct {
	UpdateID int     `json:"update_id"`
	Message  message `json:"message"`
}

type messageReply struct {
	ChatID  int
	Message string
}
