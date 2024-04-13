package main

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func getWakeUpMsg() string {
	t := time.Now()
	fmt.Printf("Day: %+v, Hour: %+v, Seconds: %+v", t.Day(), t.Hour(), t.Second())

	return ""
}

func main() {
	fmt.Println("Yo dawg whats good")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		msg := getWakeUpMsg()
		_, err := w.Write([]byte(msg))
		if err != nil {
			log.Fatal("Unable to get a formatted startup message")
		}
	})

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal("Unable to start server")
	}
}
