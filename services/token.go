package services

import (
	"context"
	"os"
	"encoding/json"
	"log"

	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

func InitializeApp() (*auth.Client, error){
	// Construct service account credentials from environment variables
	cred := map[string]string{
		"type":                        os.Getenv("FIREBASE_TYPE"),
		"project_id":                  os.Getenv("FIREBASE_PROJECT_ID"),
		"private_key_id":              os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
		"private_key":                 os.Getenv("FIREBASE_PRIVATE_KEY"),
		"client_email":                os.Getenv("FIREBASE_CLIENT_EMAIL"),
		"client_id":                   os.Getenv("FIREBASE_CLIENT_ID"),
		"auth_uri":                    os.Getenv("FIREBASE_AUTH_URI"),
		"token_uri":                   os.Getenv("FIREBASE_TOKEN_URI"),
		"auth_provider_x509_cert_url": os.Getenv("FIREBASE_AUTH_PROVIDER_X509_CERT_URL"),
		"client_x509_cert_url":        os.Getenv("FIREBASE_CLIENT_X509_CERT_URL"),
		"universe_domain":			   os.Getenv("FIREBASE_UNIVERSE_DOMAIN"),
	}
	// Convert credentials to JSON
    credJSON, err := json.Marshal(cred)
    if err != nil {
        log.Fatalf("error marshalling credentials: %v", err)
    }

	opt := option.WithCredentialsJSON(credJSON)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err!=nil {
		return nil, err
	}
	client, err := app.Auth(context.Background())
	if err!=nil{
		return nil, err
	}
	return client, nil
}

func VerifyIDToken(client *auth.Client, idToken string) (*auth.Token, error) {
    token, err := client.VerifyIDToken(context.Background(), idToken)
    if err != nil {
        return nil, err
    }
    return token, nil
}