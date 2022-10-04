package main

import (
	"log"
	"net/http"
)

func main() {
	mux := routes()

	log.Println("👀 Staring is caring 👀, PORT started on 9001")

	_ = http.ListenAndServe(":9001", mux)
}
