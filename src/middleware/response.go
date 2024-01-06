package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

type HTTPError struct {
	Error  string `json:"error"`
	Detail string `json:"detail"`
}

func NewHTTPError(w http.ResponseWriter, message string, status int, err error) {
	logMessage := message + ". err = " + err.Error() + "\n"
	log.Print(logMessage)

	resp := HTTPError{
		Error:  message,
		Detail: err.Error(),
	}
	JSONResponse(w, status, resp)
}

func JSONResponse(w http.ResponseWriter, status int, content interface{}) {
	resp, err := json.Marshal(content)
	if err != nil {
		log.Printf("Cannot encode json, err =%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
}
