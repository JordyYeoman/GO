package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Generic method to handle json responses
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
	//_, err = w.Write(data)
	//if err != nil {
	//	return
	//}
	w.Write(data)
}
