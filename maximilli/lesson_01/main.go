package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"log"
	"net/http"
)

func main() {
	// TODO:
	// 1. Create DB connection
	// 2. Add routing for basic routes
	// 3. Add auth

	host := "localhost:3000"
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	fmt.Printf("Server running on http://%s\n", host)

	router.Mount("/events", EventRoutes())

	// Start server
	err := http.ListenAndServe(host, router)
	if err != nil {
		log.Fatal("Unable to start server.")
	}

	// Defer DB disconnection
	defer func() {
		fmt.Println("Disconnecting DB + Shutting down server")
	}()
}
