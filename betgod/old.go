package main

//
//func main() {
//	fmt.Println("Full send")
//
//	// Api endpoint we want to hit
//	url := "https://jsonplaceholder.typicode.com/todos/1"
//
//	response, err := http.Get(url)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Remember, defer will run the piece of code at the end of the function scope.
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			log.Fatal(err)
//		}
//	}(response.Body)
//
//	if response.StatusCode == http.StatusOK {
//		// Second version
//		todoItem := Todo{}
//
//		decoder := json.NewDecoder(response.Body)
//		decoder.DisallowUnknownFields() // If you want to force a structure into responses from the API.
//
//		if err := decoder.Decode(&todoItem); err != nil {
//			log.Fatal("Decoder error: ", err)
//		}
//
//		fmt.Println("Decoder output")
//		fmt.Println(todoItem)
//
//		// convert Go Struct to JSON
//		todo, err := json.Marshal(todoItem)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		fmt.Println(string(todo))
//	}
//
//	return
//}
//
//func firstWayToDestructureJSON(response *http.Response) {
//	// Fist way to implement JSON destructuring
//	bodyBytes, err := io.ReadAll(response.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//	// Good way to just simply log out our response from the endpoint
//	//data :=  string(bodyByte)
//	//fmt.Println(data)
//
//	// Verbose way of unmarshalling data
//	// Create an element to place the response unmarshalled data to
//	todoItem := Todo{}
//
//	jErr := json.Unmarshal(bodyBytes, &todoItem)
//	if jErr != nil {
//		return
//	}
//
//	fmt.Printf(`Data from API: %+v`, todoItem)
//}