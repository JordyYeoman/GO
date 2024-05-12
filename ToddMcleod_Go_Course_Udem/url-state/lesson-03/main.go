package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// redirects
func main() {

	http.HandleFunc("/", foo)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/setCookie", setCookie)
	http.HandleFunc("/readCookie", readCookie)
	http.HandleFunc("/*", wildcard)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.ListenAndServe("localhost:8080", nil)
}

func foo(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Your request method at foo: ", req.Method)

	// Track amount of times user interacts / visits page
	trackAndSetCookie(w, req)
}

func trackAndSetCookie(w http.ResponseWriter, req *http.Request) {
	tCookie := "track-cookie"
	// check if cookie exists
	c, err := req.Cookie(tCookie)

	// silently handle err
	if err != nil {
		if strings.Contains(err.Error(), "named cookie not present") {
			fmt.Println("err reading cookie: ", err)
			// safe to continue
		} else {
			// handle generic cookie err
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	}

	// update the count
	if c != nil {
		newVal, err := strconv.Atoi(c.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// increase cookie value
		newVal++
		http.SetCookie(w, &http.Cookie{Name: tCookie, Value: strconv.Itoa(newVal)})
		return
	}

	// Otherwise set a new cookie
	http.SetCookie(w, &http.Cookie{Name: tCookie, Value: "1"})
}

func bar(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Your request method at bar: ", req.Method)
	http.Redirect(w, req, "/", http.StatusMovedPermanently) // 301
}

func wildcard(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Your request method at *wildcard: ", req.Method)
	//http.Redirect(w, req, "/", http.StatusTemporaryRedirect) // 307
	http.Redirect(w, req, "/", http.StatusSeeOther) // 303 - also temporary redirect
}

func setCookie(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "test-cookie", Value: "blah blah"})
	fmt.Fprintln(w, "COOKIE WRITTEN DAWG")
}

func readCookie(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("test-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, "Your cookie is: ", c)
}
