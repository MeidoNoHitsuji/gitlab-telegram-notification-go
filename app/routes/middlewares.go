package routes

import (
	"gitlab-telegram-notification-go/toggl"
	"io/ioutil"
	"net/http"
	"os"
)

func ToggleWebhookSignature(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signature := r.Header.Get("X-Webhook-Signature-256")

		if signature == "" {
			w.WriteHeader(403)
			return
		}

		secret := os.Getenv("TOGGLE_SECRET")

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		if !toggl.HmacIsValid(string(body), signature, secret) {
			w.WriteHeader(403)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
