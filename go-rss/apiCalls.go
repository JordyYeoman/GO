package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func handleApiCall() {
	url := os.Args[1]

	api(url)
}

type Pokemon struct {
	id int `json:"id"`
}

func api(url string) {
	fmt.Println(url)
	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	// read response body
	body, error := io.ReadAll(response.Body)
	if error != nil {
		fmt.Println(error)
	}

	// response

	// close response body
	response.Body.Close()

	// print response body
	fmt.Println(string(body))
}
