package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type hotdog int
type hotcat int
type GenericRes int

func (g GenericRes) ServeHTTP(res http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(res, "Generic Response on base path")
	if err != nil {
		log.Fatal("'/', Unable to write string")
	}
}

// Attaching the serve http method to the hotdog type
func (d hotdog) ServeHTTP(res http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(res, "snoop doggy dog")
	if err != nil {
		log.Fatal("'/dog', Unable to write string")
	}
}

// Attach another method to the `hotcat` type
func (c hotcat) ServeHTTP(res http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(res, "Alone in the world, the little catdawg")
	if err != nil {
		log.Fatal("'/cat', Unable to write string")
	}
}

// Cleaner alternative is to use the HandleFunc,
// so we no longer need to create a type with the handler method attached.
func e(res http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(res, "Champion Mentality")
	if err != nil {
		log.Fatal("'/e/', Unable to write string")
	}
}

func main() {
	var a GenericRes
	var d hotdog
	var c hotcat
	fmt.Println("Good morning Sir")

	mux := http.NewServeMux()
	mux.Handle("/", a)
	mux.Handle("/dog", d)
	mux.Handle("/cat", c)
	mux.HandleFunc("/e/", e)

	// Alternatively, we pass `nil` to listenAndServe, and use the
	// http.Handle / http.HandleFunc to serve our response.

	http.ListenAndServe("localhost:8080", mux)
	return
}
