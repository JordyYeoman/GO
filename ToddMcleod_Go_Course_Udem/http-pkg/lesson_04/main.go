package main

import (
	"io"
	"log"
	"net/http"
)

func g(res http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(res, "Generic Response using HandleFunc")
	if err != nil {
		log.Fatal("'/', Unable to write string")
	}
}
func main() {
	http.HandleFunc("/", g)

	http.ListenAndServe("localhost:8080", nil)
}
