package bot

import (
	"context"
	d "frigga/modules/driver"
	"time"

	"cloud.google.com/go/firestore"
)

type session struct {
	ActiveCommand string `firestore:"activeCommand"`
	UserID        string `firestore:"userId"`
	Name          string `firestore:"name"`
}

type history struct {
	Command   string    `firestore:"command"`
	Input     string    `firestore:"input"`
	Output    string    `firestore:"output"`
	Timestamp time.Time `firestore:"timestamp"`
}

// InitSession ...
func InitSession(sessionID string, name string) error {
	ctx := context.Background()
	session := session{"", sessionID, name}
	_, err := d.FS.Doc("bot_sessions/"+sessionID).Set(ctx, session)
	if err != nil {
		return err
	}
	return nil
}

// UpdateSession ...
func UpdateSession(sessionID, command string) error {
	ctx := context.Background()
	_, err := d.FS.Doc("bot_sessions/"+sessionID).Update(ctx, []firestore.Update{{Path: "activeCommand", Value: command}})
	if err != nil {
		return err
	}
	return nil
}

// GetSession ...
func GetSession(sessionID string) (string, error) {
	var command string
	var session session

	ctx := context.Background()
	docsnap, err := d.FS.Doc("bot_sessions/" + sessionID).Get(ctx)
	if err != nil {
		return command, err
	}

	if err := docsnap.DataTo(&session); err != nil {
		return command, err
	}

	return session.ActiveCommand, nil
}

// LogSession ...
func LogSession(sessionID, command, input, output string) error {
	history := history{
		Command:   command,
		Input:     input,
		Output:    output,
		Timestamp: time.Now(),
	}

	ctx := context.Background()
	_, _, err := d.FS.Collection("bot_sessions/"+sessionID+"/history").Add(ctx, history)

	return err
}
