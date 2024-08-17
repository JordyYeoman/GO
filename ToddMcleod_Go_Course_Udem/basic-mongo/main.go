package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Yo yo")
	r := mux.NewRouter()
	r.HandleFunc("/", home).Methods("GET")
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("unable to start server")
	}
}

func home(w http.ResponseWriter, res *http.Request) {
	_, err := fmt.Fprint(w, "Welcome homie!\n")
	if err != nil {
		log.Fatal("unable to load home route", err)
	}
}
