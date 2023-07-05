package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Google Gmail API helper functions
const (
//// See, edit, create, and delete all of your Google Gmail files
//GmailScope = "https://www.googleapis.com/auth/gmail"
//
//// See, edit, create, and delete only the specific Google Gmail files
//// you use with this app
//GmailFileScope = "https://www.googleapis.com/auth/gmail.file"
//
//// See and download all your Google Gmail files
//GmailReadonlyScope = "https://www.googleapis.com/auth/gmail.readonly"
)

// Use this as a singleton object for the Google Sheets service
var gmailService *gmail.Service

// Retrieve a token, saves the token, then returns the generated client.
func getGmailClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "secrets/gmail_token.json"
	tok, err := tokenGmailFromFile(tokFile)
	if err != nil {
		tok = getGmailTokenFromWeb(config)
		saveGmailToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getGmailTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenGmailFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveGmailToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func Gmail() (*gmail.Service, error) {
	if gmailService != nil {
		return gmailService, nil
	}

	ctx := context.Background()
	b, err := ioutil.ReadFile("secrets/client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailMetadataScope, gmail.GmailSendScope, gmail.GmailModifyScope, gmail.GmailReadonlyScope) //drive.DriveMetadataReadonlyScope
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getGmailClient(config)

	newService, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err == nil {
		gmailService = newService
	}
	return newService, err

}
