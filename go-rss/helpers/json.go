package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// marshall the payload into a JSON string
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshall json response: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200) // setting response status code
	w.Write(data)
}
