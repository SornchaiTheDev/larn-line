package services

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func NewFirestore() (*firestore.Client, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("./serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	firestore, err := app.Firestore(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return firestore, nil
}
