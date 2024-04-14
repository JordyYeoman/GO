package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func BonusRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("This is the base '/' path for bonus routes!"))
		if err != nil {
			log.Fatal("Error occurred")
		}
	})

	return r
}
