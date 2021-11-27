package api

import (
	"github.com/bitly/go-simplejson"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func RespondDeepLinkWithJSON(w http.ResponseWriter, code int, converterResponse string) {
	json := simplejson.New()
	json.Set("deeplink", converterResponse)
	payload, _ := json.MarshalJSON()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(payload)
}
func RespondError(w http.ResponseWriter, code int, message string) {
	e := ErrorResponse{
		Code:    code,
		Status:  "Error",
		Message: message,
	}
	json := simplejson.New()
	json.Set("error", e)
	a, _ := json.MarshalJSON()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(a)
}
func RespondWebURLWithJSON(w http.ResponseWriter, code int, converterResponse string) {
	json := simplejson.New()
	json.Set("weburl",converterResponse)
	payload, _ := json.MarshalJSON()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(payload)
}