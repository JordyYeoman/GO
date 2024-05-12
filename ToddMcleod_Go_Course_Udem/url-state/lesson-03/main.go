package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// redirects
func main() {

	http.HandleFunc("/", foo)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/setCookie", setCookie)
	http.HandleFunc("/readCookie", readCookie)
	http.HandleFunc("/createSession", createSession)
	http.HandleFunc("/expireCookie", expireCookie)
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

	if err == http.ErrNoCookie {
		http.SetCookie(w, &http.Cookie{Name: tCookie, Value: "1"})
		return
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
		io.WriteString(w, "# of times visiting site: "+strconv.Itoa(newVal))
		return
	}
}

func createSession(w http.ResponseWriter, req *http.Request) {
	fmt.Println("TIme: ", time.Now().Format(""))
	http.SetCookie(w, &http.Cookie{Name: "session", Value: "Now dude!"})
	//fmt.Fprintln(w, "Session started")
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

func expireCookie(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session")
	if err != nil {
		http.Redirect(w, req, "/set", http.StatusSeeOther)
		return
	}

	c.MaxAge = -1
	http.SetCookie(w, c)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}
