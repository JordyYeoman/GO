package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Todo struct {
	UserID    int    `json:"userId"` // Gives us a way to map the JSON values to go struct properties
	ID        int    `json:"id"`
	Title     string `json:"title"` // Uppercase means the value is 'exported'
	Completed bool   `json:"completed"`
}

func main() {
	fmt.Println("Full send")

	// Api endpoint we want to hit
	url := "https://jsonplaceholder.typicode.com/todos/1"

	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	// Remember, defer will run the piece of code at the end of the function scope.
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(response.Body)

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		// Good way to just simply log out our response from the endpoint
		//data :=  string(bodyByte)
		//fmt.Println(data)

		// Verbose way of unmarshalling data
		// Create an element to place the response unmarshalled data to
		todoItem := Todo{}

		jErr := json.Unmarshal(bodyBytes, &todoItem)
		if jErr != nil {
			return
		}

		fmt.Printf(`Data from API: %+v`, todoItem)
	}

	return
}
