package toggl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type MethodType string

const (
	ApiType     MethodType = "api/v9"
	WebhookType MethodType = "webhooks/api/v1"
)

func builder(typeUrl MethodType, method string, url string, token string, data []byte) (string, error) {
	toggleDomain := os.Getenv("TOGGLE_DOMAIN")

	var (
		req *http.Request
		err error
	)

	if data == nil {
		req, err = http.NewRequest(method,
			fmt.Sprintf("%s/%s/%s", toggleDomain, typeUrl, url), nil)
	} else {
		req, err = http.NewRequest(method,
			fmt.Sprintf("%s/%s/%s", toggleDomain, typeUrl, url), bytes.NewBuffer(data))
	}

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	if token != "" {
		req.SetBasicAuth(token, "api_token")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(string(body))
	}

	return string(body), nil
}

func Me(token string) (*UserType, error) {
	response, err := builder(ApiType, http.MethodGet, "me", token, nil)

	if err != nil {
		return nil, err
	}

	var result UserType

	err = json.Unmarshal([]byte(response), &result)

	if err != nil {
		return nil, err
	} else {
		return &result, nil
	}
}

func Events() (*[]Event, error) {
	response, err := builder(WebhookType, http.MethodGet, "event_filters", "", nil)

	if err != nil {
		return nil, err
	}

	var responseData map[string][]string

	err = json.Unmarshal([]byte(response), &responseData)

	if err != nil {
		return nil, err
	} else {
		var result []Event

		for key, entities := range responseData {
			event := Event{
				Action:   key,
				Entities: entities,
			}

			result = append(result, event)
		}

		return &result, nil
	}
}

func CreateSubscriptions(workspaceId int, token string, data SubscriptionData) (*SubscriptionData, error) {
	data.Secret = os.Getenv("TOGGLE_SECRET")

	out, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	response, err := builder(WebhookType, http.MethodPost, fmt.Sprintf("subscriptions/%d", workspaceId), token, out)

	if err != nil {
		return nil, err
	}

	var result SubscriptionData

	err = json.Unmarshal([]byte(response), &result)

	if err != nil {
		return nil, err
	} else {
		return &result, nil
	}
}

func UpdateSubscriptions(workspaceId int, subscriptionId int, token string, data SubscriptionData) (*SubscriptionData, error) {
	data.Secret = os.Getenv("TOGGLE_SECRET")

	out, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	response, err := builder(WebhookType, http.MethodPut, fmt.Sprintf("subscriptions/%d/%d", workspaceId, subscriptionId), token, out)

	if err != nil {
		return nil, err
	}

	var result SubscriptionData

	err = json.Unmarshal([]byte(response), &result)

	if err != nil {
		return nil, err
	} else {
		return &result, nil
	}
}

func GetSubscriptions(workspaceId int, token string) ([]*SubscriptionData, error) {
	response, err := builder(WebhookType, http.MethodGet, fmt.Sprintf("subscriptions/%d", workspaceId), token, nil)

	if err != nil {
		return nil, err
	}

	var result []*SubscriptionData

	err = json.Unmarshal([]byte(response), &result)

	if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}
