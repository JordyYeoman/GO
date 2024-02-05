package main

import (
	"database/sql"
	"fmt"
	"log"
)

func handleSQLExampleQueries(db *sql.DB) {
	// 1. Testing table creation
	createProductTable(db)

	// 2. Testing product creation
	product := Product{"Book", 12.33, true}
	pk := insertProduct(db, product)
	fmt.Printf("Price key: %v\n", pk)

	// 3. Testing a basic query
	var name string
	var price float64
	var available bool
	// invalidRow := 999999 // testing error case

	query := "SELECT name, price, available FROM product WHERE id = ?"
	// queryErr := db.QueryRow(query, pk).Scan(&name, &price, &available)    // Valid query
	queryErr := db.QueryRow(query, pk).Scan(&name, &price, &available) // No rows to be found
	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			log.Fatalf("No rows found for the id: %d", pk) // Handle logic for no rows being found
		}
		fmt.Printf("Error: %e", queryErr)
		log.Fatal(queryErr)
	}
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Price: %f\n", price)
	fmt.Printf("Name: %t\n", available)

	// 4. Multiple rows query
	data := []Product{}
	rows, err := db.Query("SELECT name, available, price from product")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	// to SCAN db vals
	// var name string
	// var available bool
	// var price float64

	for rows.Next() {
		err := rows.Scan(&name, &available, &price)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Product{name, price, available})
	}

	fmt.Println(data)
}
