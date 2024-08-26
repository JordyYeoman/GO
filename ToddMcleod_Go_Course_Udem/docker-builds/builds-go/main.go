package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("unable to start server")
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Hello from test go build - inside a docker container")
	if err != nil {
		log.Println("unable to write string")
	}
}
