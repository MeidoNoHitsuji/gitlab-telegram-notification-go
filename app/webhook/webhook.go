package webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xanzy/go-gitlab"
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
func (hook Webhook) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	event, err := hook.parse(request)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte(fmt.Sprintf("could parse the webhook event: %v", err)))
		return
	}

	// Handle the event before we return.
	if err := hook.handle(event); err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte(fmt.Sprintf("error handling the event: %v", err)))
		return
	}

	// Write a response when were done.
	writer.WriteHeader(204)
}

func (hook Webhook) handle(event interface{}) error {
	str, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("could not marshal json event for logging: %v", err)
	}

	// Just write the event for this example.
	fmt.Println(string(str))

	return nil
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
