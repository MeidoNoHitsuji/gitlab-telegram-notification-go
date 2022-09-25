package routes

import (
	"fmt"
	"net/http"
)

func ToggleWebhookSignature(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signature := r.Header.Get("X-Webhook-Signature-256")

		if signature == "" {
			w.WriteHeader(403)
			return
		}

		//secret := os.Getenv("TOGGLE_SECRET")

		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(403)
			fmt.Println(err)
			return
		}

		fmt.Println(r.Form)
		fmt.Println(r.PostForm)

		//if !toggl.HmacIsValid(string(bodyData), signature, secret) {
		//	w.WriteHeader(403)
		//	return
		//}

		next.ServeHTTP(w, r)
	})
}
