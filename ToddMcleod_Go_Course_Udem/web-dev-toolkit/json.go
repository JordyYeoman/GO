package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting server")
	http.HandleFunc("/", boo)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("unable to start server")
	}
}

func boo(w http.ResponseWriter, r *http.Request) {

}
