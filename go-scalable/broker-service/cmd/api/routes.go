package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func routes() http.Handler {
	// mux is the common name used when routing in go
	// mux := chi.NewRouter
	mux := chi.NewRouter()

	// Specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://*",
			"http://*",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// return
}
