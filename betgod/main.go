package main

import (
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// TODO: Add weighting from betfair model: https://www.betfair.com.au/hub/sports/afl/afl-predictions-model/
// TODO: Add endpoints for analysis between teams, quarter, win prob % etc
// TODO: Handle simple auth

func main() {
	// Connect to DB
	db := connectToDB()

	// Create routes for API
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

	router.Mount("/teams", TeamRoutes(db))

	// Spin up server
	http.ListenAndServe(":3000", router)

	// Disconnect DB
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.WithError(err).Warn("Failed to disconnect DB")
		}
	}(db) // Defer means run this when the wrapping function terminates
}
