package main

import (
	"database/sql"
	"log"
)

type Product struct {
	Name      string
	Price     float64
	Available bool
}

func createProductTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS product (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price NUMERIC(6,2) NOT NULL,
		available BOOLEAN,
		created timestamp DEFAULT NOW()
	)`

	_, err := db.Exec(query) // Execute query against DB without returning any rows
	if err != nil {
		log.Fatal(err)
	}
}

// func insertProduct(db *sql.DB, product Product) int {
// 	query := `INSERT INTO product (name, price, available)
// 			  VALUES ($1, $2, $3)`

// 	var pk int
// 	err := db.QueryRow(query, product.Name, product.Price, product.Available).Scan(&pk)

// 	fmt.Printf("PK value: %v", pk)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return pk
// }

func insertProduct(db *sql.DB, product Product) int {
	query := "INSERT INTO product (name, price, available) VALUES (?, ?, ?);"
	result, err := db.Exec(query, product.Name, product.Price, product.Available)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the last inserted ID
	pk, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return int(pk)
}
