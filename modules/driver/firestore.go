package driver

import (
	"context"
	"log"
	"os"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

// FS ...
var FS *firestore.Client

// FSContext ...
// var FSContext *context.Context

// InitializeFirestore ...
func InitializeFirestore() {
	ctx := context.Background()
	sa := option.WithCredentialsFile(os.Getenv("FIREBASE_ACCOUNT_URL"))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	FS = client
	// FSContext = &ctx
}
