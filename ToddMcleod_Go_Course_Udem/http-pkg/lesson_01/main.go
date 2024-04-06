package main

import (
	"fmt"
	"net/http"
)

// Important!!
// type handler interface {
//	ServeHttp(ResponseWriter, Request)
// }

type hotdog int

// Extending the above hotdog type to include
func (m hotdog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Any code you want in this cheeky bugger")
}

func main() {
	fmt.Println("Yo dawg whats up!!")

	var d hotdog

	http.ListenAndServe(":8080", d)

	return
}
