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

func api(url string) {
	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	// read response body
	body, error := io.ReadAll(response.Body)
	if error != nil {
		fmt.Println(error)
	}

	// Do something with the body response.
	// response.Body.Read(p)

	// close response body
	response.Body.Close()

	// print response body
	fmt.Println(string(body))
}
