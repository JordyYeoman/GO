package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Unique Identifiers")

	http.HandleFunc("/", foo)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	err := http.ListenAndServe("localhost:8080", nil) // pass nil so http package uses default mux router
	if err != nil {
		return
	}
}

func foo(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session-id")
	if err != nil {
		id, err := uuid.NewUUID()
		if err != nil {
			log.Fatal("Unable to create UUID")
		}
		cookie = &http.Cookie{
			Name:  "session-id",
			Value: id.String(),
			// Secure: true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}

	fmt.Println(cookie)
}
