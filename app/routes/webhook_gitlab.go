package routes

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/gitclient"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Webhook struct {
	Secret         string
	EventsToAccept []gitlab.EventType
}

// ServeHTTP tries to parse Gitlab events sent and calls handle function
// with the successfully parsed events.
func (hook Webhook) ServeHTTP(c *gin.Context) {
	event, err := hook.parse(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Could parse the webhook event: %v", err),
		})
		return
	}

	// Handle the event before we return.
	if err := gitclient.Handler(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Error handling the event: %v", err),
		})
		return
	}

	// Write a response when were done.
	c.Status(http.StatusNoContent)
}

// parse verifies and parses the events specified in the request and
// returns the parsed event or an error.
func (hook Webhook) parse(r *http.Request) (interface{}, error) {
	defer func() {
		if _, err := io.Copy(ioutil.Discard, r.Body); err != nil {
			log.Printf("could discard request body: %v", err)
		}
		if err := r.Body.Close(); err != nil {
			log.Printf("could not close request body: %v", err)
		}
	}()

	if r.Method != http.MethodPost {
		return nil, errors.New("invalid HTTP Method")
	}

	// If we have a secret set, we should check if the request matches it.
	if len(hook.Secret) > 0 {
		signature := r.Header.Get("X-Gitlab-Token")
		if signature != hook.Secret {
			return nil, errors.New("token validation failed")
		}
	}

	event := r.Header.Get("X-Gitlab-Event")
	if strings.TrimSpace(event) == "" {
		return nil, errors.New("missing X-Gitlab-Event Header")
	}

	eventType := gitlab.EventType(event)
	if !isEventSubscribed(eventType, hook.EventsToAccept) {
		return nil, errors.New("event not defined to be parsed")
	}

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil || len(payload) == 0 {
		return nil, errors.New("error reading request body")
	}

	return gitlab.ParseWebhook(eventType, payload)
}

func isEventSubscribed(event gitlab.EventType, events []gitlab.EventType) bool {
	for _, e := range events {
		if event == e {
			return true
		}
	}
	return false
}
