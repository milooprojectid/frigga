package repository

import (
	"context"
	d "frigga/modules/driver"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type session struct {
	ActiveCommand string `firestore:"activeCommand"`
	Provider      string `firestore:"provider"`
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
func InitSession(provider string, sessionID string, name string) error {
	ctx := context.Background()
	session := session{"", provider, sessionID, name}
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

// GetCovid19Data ...
func GetCovid19Data() (Covid19Data, error) {
	var covid19Data Covid19Data

	ctx := context.Background()
	docsnap, err := d.FS.Doc("bot_data/covid19").Get(ctx)
	if err != nil {
		return covid19Data, err
	}

	if err := docsnap.DataTo(&covid19Data); err != nil {
		return covid19Data, err
	}

	return covid19Data, nil
}

// SetCovid19SubsData ...
func SetCovid19SubsData(uid string, provider string) error {
	ctx := context.Background()
	payload := map[string]interface{}{
		"is_subscribe": true,
		"provider":     provider,
		"token":        uid,
	}
	_, err := d.FS.Doc("bot_data/covid19/subscriptions/"+uid).Set(ctx, payload)

	return err
}

// GetCovid19Subscription ...
func GetCovid19Subscription() ([]Covid19Subcription, error) {
	covid19Subscriptions := []Covid19Subcription{}
	var sub Covid19Subcription

	ctx := context.Background()
	iter := d.FS.Collection("bot_data/covid19/subscriptions").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return covid19Subscriptions, err
		}
		doc.DataTo(&sub)
		covid19Subscriptions = append(covid19Subscriptions, sub)
	}

	return covid19Subscriptions, nil
}

// GetAllBotSessions ...
func GetAllBotSessions(isProd bool) ([]BotSession, error) {
	botSessions := []BotSession{}
	var session BotSession
	var path string

	if isProd {
		path = "bot_sessions"
	} else {
		path = "bot_sessions_test"
	}

	ctx := context.Background()
	iter := d.FS.Collection(path).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return botSessions, err
		}
		doc.DataTo(&session)
		botSessions = append(botSessions, session)
	}

	return botSessions, nil
}
