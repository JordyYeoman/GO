package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5xx Error %v", msg)
	}

	type errResponse struct {
		Error string `json:"Error"`
	}
	// The struct above ^^ will marshal into a json object, eg:
	// {
	// 	"error": "Something went wrong"
	// }

	respondWithJSON(w, code, errResponse{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// marshall the payload into a JSON string
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshall json response: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code) // setting response status code
	w.Write(data)
}
