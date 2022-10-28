package main

import (
	"database/sql"
	"log"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service, SEND IT")

	// TODO connect to DB

	// setup config
	app := Config{}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%f", webPort)
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
