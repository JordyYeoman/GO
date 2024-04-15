package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// Practice Q's
func HomeRoute(res http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	_, err := io.WriteString(res, "This is the handle route using our type Handler")
	if err != nil {
		log.Fatal("'/', Unable to write string")
	}
}

func DogeRoute(res http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(res, "This is our DOGE ROUTE")
	if err != nil {
		log.Fatal("/Dogeroute error")
	}
}

func main() {
	fmt.Println("Server starting")

	mux := http.NewServeMux()
	mux.HandleFunc("/", HomeRoute)
	mux.HandleFunc("/doge/", DogeRoute)

	http.ListenAndServe("localhost:8080", mux)
	return
}
