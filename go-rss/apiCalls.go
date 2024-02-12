package main

import (
	"encoding/json"
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
	Id int `json:"id"` // Struct values need to be exported by having a capital I.
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

	var p Pokemon
	// parse the json data
	pErr := json.Unmarshal(body, &p)
	if pErr != nil {
		fmt.Println(pErr)
		return
	}

	fmt.Println(p)

	// close response body
	response.Body.Close()

	// print response body
	// fmt.Println(string(body))
}
