package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Server Starting...")

	h1 := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World\n")
		io.WriteString(w, r.Method)
	}
	http.HandleFunc("/", h1)

	fmt.Println("Server running at Port:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
