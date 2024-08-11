package main

import (
	"database/sql"
	"fmt"
	"github.com/JordyYeoman/GO/ToddMcleod_Go_Course_Udem/test-deployments/lesson-02/config"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func connectDB() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", config.Envs.DbUsername+":"+config.Envs.DbPassword+"@tcp("+config.Envs.DbConnectionUrl+":"+config.Envs.DbPort+")/"+config.Envs.DbName+"?charset=utf8")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	fmt.Println("Yo welcome brotherrrrr")

	db, err := connectDB()
	if err != nil {
		log.Fatal("Unable to connect to DB", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Unable to close DB connection")
		}
	}(db)

	// If we have a db connection, test it
	err = db.Ping()
	if err != nil {
		log.Fatal("Unable to ping db", err)
	}
}
