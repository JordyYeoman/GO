package main

import (
	"database/sql"
	"log"
)

const webPort = "80"
var counts int64

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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := osGetenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready....")
			counts++
		} else {
			log.Println("Connected to Postgres database!")
			return connection
		}

		
	}
}
