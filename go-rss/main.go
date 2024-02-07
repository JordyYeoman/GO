package main

import (
	// "database/sql"
	"fmt"
	// "log"
	// "net/http"
	// "os"
	// "github.com/go-chi/chi"
	// "github.com/go-chi/cors"
	// _ "github.com/go-sql-driver/mysql" // Importing a package for side effects, no direct usages (interface for DB)
	// "github.com/joho/godotenv"
)

func newCard() string {
	return "Five of diamonds"
}

func main() {
	fmt.Println("Hello World")
	fmt.Println("Server starting...")

	card := newCard()
	fmt.Println(card)

	cards := []string{newCard(), newCard()}
	cards = append(cards, "Ace of Spades")

	for i, card := range cards {
		fmt.Println(i, card)
	}

	handleApiCall()

	// Array - fixed list of things
	// Slice - An array that can grow or shrink (dynamic arr)
	// Every element of a slice must be of identical type

	// godotenv.Load()

	// port := os.Getenv("PORT")
	// db_url := os.Getenv("DB_URL")

	// fmt.Println("Connecting to DB:")
	// db, dbErr := sql.Open("mysql", db_url)
	// if dbErr != nil {
	// 	log.Fatal(dbErr)
	// }

	// if dbErr = db.Ping(); dbErr != nil {
	// 	log.Fatal(dbErr)
	// }

	// // Working with SQL examples
	// handleSQLExampleQueries(db)

	// defer db.Close() // Defer means run this when the wrapping function terminates

	// if port == "" {
	// 	log.Fatal("PORT is not found in the environment")
	// }

	// router := chi.NewRouter()

	// router.Use(cors.Handler(cors.Options{
	// 	// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
	// 	AllowedOrigins: []string{"https://*", "http://*"},
	// 	// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	// 	ExposedHeaders:   []string{"Link"},
	// 	AllowCredentials: false,
	// 	MaxAge:           300, // Maximum value not ignored by any of major browsers
	// }))

	// v1Router := chi.NewRouter()

	// // v1Router.HandleFunc("/healthz", handlerReadiness) // Any endpoint scoping
	// v1Router.Get("/healthz", handlerReadiness) // Scopes the route to only GET requests
	// v1Router.Get("/err", handlerError)

	// router.Mount("/v1", v1Router)

	// srv := &http.Server{
	// 	Handler: router,
	// 	Addr:    ":" + port,
	// }

	// log.Printf("Server starting on PORT: %v", port)
	// err := srv.ListenAndServe()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("PORT:", port)
}
