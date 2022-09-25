package routes

import (
	"fmt"
	"io/ioutil"
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

		body := r.Body
		bodyData, err := ioutil.ReadAll(body)

		if err != nil {
			fmt.Println("ReadAll")
			fmt.Println(err)
			w.WriteHeader(403)
			return
		}

		body.Close()

		fmt.Println(bodyData)

		err = r.ParseForm()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(403)
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
