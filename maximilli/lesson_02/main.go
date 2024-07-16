package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	fmt.Println("Server running...")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Error occurred: ", err)
		return
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	writeString, err := io.WriteString(w, "All good")
	if err != nil {
		log.Fatal("Error occurred: ", writeString)
		return
	}

	host := "localhost:3000"
	router := chi.NewRouter()

	fmt.Printf("Server running on http://%s\n", host)

	// Start server
	err = http.ListenAndServe(host, router)
}
