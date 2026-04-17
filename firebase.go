package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var app *firebase.App
var firestoreClient *firestore.Client

func initFirebase() {
	ctx := context.Background()

	opt := option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON")))

	var err error

	// 🔥 Initialize Firebase app
	app, err = firebase.NewApp(ctx, &firebase.Config{
	ProjectID: "event-backend-e6450",
}, opt)
	if err != nil {
		log.Fatal("Firebase init failed:", err)
	}

	// 🔥 Initialize Firestore (ONLY ONCE)
	firestoreClient, err = app.Firestore(ctx)
	if err != nil {
		log.Fatal("Firestore init failed:", err)
	}

	log.Println("Firebase connected ✅")
}