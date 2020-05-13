package googlesheet

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// Message - Structure of feedback
type Message struct {
	time     string
	context  string
	Feedback string
}

//NewFeedback creates a new in googlesheet feedback
func NewFeedback(c string, t string) *Message {
	return &Message{
		time:     time.Now().String(),
		context:  c,
		Feedback: t,
	}
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		fmt.Printf("Unable to read authorization code: %v", err)

	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		fmt.Printf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
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
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Send - Init func
func (f Message) Send(spreadsheetID string) (string, bool) {
	b, err := ioutil.ReadFile("config/credentials.json")
	if err != nil {
		return fmt.Sprintf("Unable to read client secret file: %v", err), false
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return fmt.Sprintf("Unable to parse client secret file to config: %v", err), false
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		return fmt.Sprintf("Unable to retrieve Sheets client: %v", err), false
	}
	Range := "A1"

	var vr sheets.ValueRange

	myval := []interface{}{f.time, f.context, f.Feedback}
	vr.Values = append(vr.Values, myval)

	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, Range, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Sprintf("Unable to retrieve data from sheet: %v", err), false
	}
	return "Sent", true
}
