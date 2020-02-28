package repository

import (
	"context"
	d "frigga/driver"
	"time"

	"cloud.google.com/go/firestore"
)

type session struct {
	Command   string `firestore:"command"`
	FirstName string `firestore:"firstName"`
	LastName  string `firestore:"lastName"`
}

type history struct {
	Command   string `firestore:"command"`
	Payload   string `firestore:"payload"`
	Output    string `firestore:"output"`
	Timestamp int64  `firestore:"timestamp"`
}

// InitSession ...
func InitSession(sessionID, firstName, lastName string) error {
	ctx := context.Background()
	session := session{"", firstName, lastName}
	_, err := d.FS.Doc("bots/telegram/sessions/"+sessionID).Set(ctx, session)
	if err != nil {
		return err
	}
	return nil
}

// UpdateSession ...
func UpdateSession(sessionID, command string) error {
	ctx := context.Background()

	_, err := d.FS.Doc("bots/telegram/sessions/"+sessionID).Update(ctx, []firestore.Update{{Path: "command", Value: command}})
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
	docsnap, err := d.FS.Doc("bots/telegram/sessions/" + sessionID).Get(ctx)
	if err != nil {
		return command, err
	}

	if err := docsnap.DataTo(&session); err != nil {
		return command, err
	}

	return session.Command, nil
}

// LogHistory ...
func LogHistory(sessionID, command, payload, output string) error {
	history := history{
		Command:   command,
		Payload:   payload,
		Output:    output,
		Timestamp: time.Now().Unix(),
	}

	ctx := context.Background()
	_, _, err := d.FS.Collection("bots/telegram/sessions/"+sessionID+"/history").Add(ctx, history)

	return err
}
