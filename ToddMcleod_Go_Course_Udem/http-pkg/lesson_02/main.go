package main

import (
	"fmt"
	"net/http"
)

type GenericStruct int

func (m GenericStruct) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Yeoman-Key", "This is a yeoman key")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, "<h1>Yo dawwggg whats up!</h1>")
}

func main() {
	fmt.Println("Basic Headers")

	var d GenericStruct
	http.ListenAndServe(":8089", d)
	return
}
