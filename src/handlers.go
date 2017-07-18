package newsletter

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"google.golang.org/appengine"
)

type SendRequest struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	v := r.URL.Query()
	email := v.Get("email")

	if err := subscribe(ctx, email); err == nil {
		responseJSON(w, "ok")
	} else {
		responseError(w, err.Error(), http.StatusBadRequest)
	}
}

func sendMailHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var sr SendRequest
	err = json.Unmarshal(data, &sr)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := appengine.NewContext(r)
	if err := send(ctx, sr.Subject, sr.Message); err == nil {
		responseJSON(w, "ok")
	} else {
		responseError(w, err.Error(), http.StatusInternalServerError)
	}
}