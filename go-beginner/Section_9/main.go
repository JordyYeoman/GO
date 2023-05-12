package main

import (
	"fmt"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	// with FPrint, instead of printing to the standard out, we print to 
	// the response writer, which in the below example is 'w'
	fmt.Fprint(w, "Hello World")
}

// The first 1024 ports on a computer are protected, can only run on port 80 in prod
func main() {
	port := "8080"
	http.HandleFunc("/", homePage)

	// Start production ready server on port 8080
	log.Println("Starting web server on PORT: ", port)
	http.ListenAndServe(":" + port, nil)
}