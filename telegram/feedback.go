package telegram

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Message - Structure of feedback
type Message struct {
	time     string
	context  string
	Feedback string
	ChatID   string
}

//NewFeedback creates a new in googlesheet feedback
func NewFeedback(c string, t string, chat string) *Message {
	return &Message{
		time:     time.Now().String(),
		context:  c,
		Feedback: t,
		ChatID:   chat,
	}
}

// Send - Init func
func (f Message) Send(BotToken string) (string, bool) {
	var telegramAPI string = "https://api.telegram.org/bot" + BotToken + "/sendMessage"
	response, err := http.PostForm(
		telegramAPI,
		url.Values{
			"chat_id": {f.ChatID},
			"text":    {f.Feedback},
		})

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return err.Error(), false
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return errRead.Error(), false
	}

	bodyString := string(bodyBytes)
	return bodyString, true
}
