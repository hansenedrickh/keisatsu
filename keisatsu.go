package keisatsu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

const (
	PanicLevel = "panic"
	ErrorLevel = "error"
	InfoLevel  = "info"
)

type Service interface {
	WatchPanic()
	Error(msg string)
	Info(msg string)
	sendWebhook(log Message)
}

// Keisatsu struct is the struct which will be used for initiating a keisatsu package
type Keisatsu struct {
	AppName     string
	WebhookURL  string
	SecretToken string
}

// Message will be the message to be sent to the webhook in case panic is happening or error given by the app
type Message struct {
	AppName    string `json:"app_name"`
	Level      string `json:"level"`
	Message    string `json:"message"`
	StackTrace string `json:"stack_trace"`
}

// New is for initiating a new Keisatsu client
func New(appName, webhookURL, secretToken string) Service {
	return &Keisatsu{
		AppName:     appName,
		WebhookURL:  webhookURL,
		SecretToken: secretToken,
	}
}

// WatchPanic is the function which you can register anywhere to avoid panic and send the panic to the webhook URL
func (k Keisatsu) WatchPanic() {
	if msg := recover(); msg != nil {
		message := Message{
			AppName:    k.AppName,
			Level:      PanicLevel,
			Message:    fmt.Sprintf("%v", msg),
			StackTrace: string(debug.Stack()),
		}
		k.sendWebhook(message)
	}
}

// Message is the function which you can use to send an error level webhook
func (k Keisatsu) Error(msg string) {
	message := Message{
		AppName: k.AppName,
		Level:   ErrorLevel,
		Message: msg,
	}
	k.sendWebhook(message)
}

// Info is the function which you can use to send an info level webhook
func (k Keisatsu) Info(msg string) {
	message := Message{
		AppName: k.AppName,
		Level:   InfoLevel,
		Message: msg,
	}
	k.sendWebhook(message)
}

func (k Keisatsu) sendWebhook(m Message) {
	requestBody, err := json.Marshal(m)
	if err != nil {
		log.Println("Failed to marshal: ", err)
	}

	req, err := http.NewRequest("POST", k.WebhookURL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error creating a request: ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Secret-Token", k.SecretToken)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println("Error sending a request: ", err)
	}
}
