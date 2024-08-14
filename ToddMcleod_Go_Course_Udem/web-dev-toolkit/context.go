package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", foo)
	http.HandleFunc("/mshl", mshl)
	http.HandleFunc("/encd", encd)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("unable to start server")
	}
}

func foo(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log.Println(ctx)
	_, err := fmt.Fprintln(w, ctx)
	if err != nil {
		log.Fatal("unable to print context")
	}
}
