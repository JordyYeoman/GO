package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type NewMarvelMoviePayload struct {
	DaysUntil           int `json:"days_until"`
	FollowingProduction struct {
		DaysUntil   int    `json:"days_until"`
		Id          int    `json:"id"`
		Overview    string `json:"overview"`
		PosterUrl   string `json:"poster_url"`
		ReleaseDate string `json:"release_date"`
		Title       string `json:"title"`
		Type        string `json:"type"`
	} `json:"following_production"`
	Id          int    `json:"id"`
	Overview    string `json:"overview"`
	PosterUrl   string `json:"poster_url"`
	ReleaseDate string `json:"release_date"`
	Title       string `json:"title"`
	Type        string `json:"type"`
}

func main() {
	url := "https://www.whenisthenextmcufilm.com/api"

	var newMarvelMovie NewMarvelMoviePayload

	// New http request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		_ = fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/json")

	// Create a client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		_ = fmt.Errorf("error sending request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("unable to close the body of response")
		}
	}(resp.Body)

	// Unmarshal data
	body, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &newMarvelMovie)
	if err != nil {
		log.Fatal("unable to unmarshal response body")
	}

	fmt.Println("Payload: ", newMarvelMovie.FollowingProduction)
}
