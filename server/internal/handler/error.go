package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func handleError(
	w http.ResponseWriter,
	err error,
	statusCode int,
	shouldLog bool,
) {
	if shouldLog {
		log.Println(err.Error())
	}

	errorBody, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{
		Error: err.Error(), // Change to something generic if this API is publicly exposed
	})

	w.WriteHeader(statusCode)
	w.Write(errorBody)
}
