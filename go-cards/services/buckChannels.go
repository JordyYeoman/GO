package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Checking websites")
	links := []string{
		"https://google.com",
		"https://facebook.com",
		"https://stackoverflow.com",
		"https://golang.org",
		"https://amazon.com",
	}

	c := make(chan string)

	for _, link := range links {
		go checkLink(link, c)
	}

	for l := range c {
		go func(link string) {
			time.Sleep(5 * time.Second)
			checkLink(link, c)
		}(l)
	}
}

func checkLink(url string, c chan string) {
	_, err := http.Get(url)
	if err != nil {
		fmt.Println("Link "+url, " might be down!")
		c <- url
		return
	}

	fmt.Println("Link "+url, "is OK")
	c <- url
}
