package firebaseadmin

import (
	"context"
	"donation-mgmt/src/config"
	"donation-mgmt/src/libs/logger"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var appClient *firebase.App
var authClient *auth.Client

func Bootstrap(appConfig *config.AppConfiguration) {
	l := logger.ForComponent("firebase-admin")
	clientOptions := []option.ClientOption{}

	if appConfig.GCPServiceAccountJSONPath != "" {
		l.Info("Initializating Firebase Admin with service account JSON")
		clientOptions = append(clientOptions, option.WithCredentialsFile(appConfig.GCPServiceAccountJSONPath))
	}

	var err error
	appClient, err = firebase.NewApp(context.Background(), nil, clientOptions...)

	if err != nil {
		panic("Error initializing Firebase client: " + err.Error())
	}

	authClient, err = appClient.Auth(context.Background())
	if err != nil {
		panic("Error initializing Firebase Auth client: " + err.Error())
	}
}

func AuthClient() *auth.Client {
	if authClient == nil {
		panic("Firebase Admin was not bootstrapped")
	}

	return authClient
}
