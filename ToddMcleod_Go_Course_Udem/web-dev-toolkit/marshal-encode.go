package main

import (
	json2 "encoding/json"
	"log"
	"net/http"
)

// Person Remember when writing to JSON, the struct property fields need to be uppercase,
// otherwise they will not get written to the response
type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func main() {
	http.Handle("/", http.NotFoundHandler())
	http.HandleFunc("/mshl", mshl)
	http.HandleFunc("/encd", encd)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("unable to start server")
	}
}

func mshl(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p1 := Person{
		Age:       29,
		FirstName: "Jordy",
		LastName:  "Yeoman",
	}

	json, err := json2.Marshal(p1)
	if err != nil {
		log.Println(err)
	}

	_, err = w.Write(json)
	if err != nil {
		log.Println("unable to write json")
	}
}

func encd(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p1 := Person{
		Age:       29,
		FirstName: "Jordy",
		LastName:  "Yeoman",
	}
	// Encoder needs to know the stream of where it will write to
	// IE - the path from which the request came from, that we can stream our json back as a response.
	err := json2.NewEncoder(w).Encode(p1)
	if err != nil {
		log.Println(err)
	}
}
