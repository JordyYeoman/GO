package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type User struct {
	UserName string
	Password string
	First    string
	Last     string
	Age      int
}

var tpl *template.Template
var dbUsers = map[string]User{}      // user ID, user
var dbSessions = map[string]string{} // session ID, user ID

func init() {
	fmt.Println("Initializing Templates")
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	fmt.Println("User Unique Session ID + lookup")

	http.HandleFunc("/", index)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/signup", signup)
	http.Handle("/favivon.ico", http.NotFoundHandler())
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("Unable to start server")
		return
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)

	err := tpl.ExecuteTemplate(w, "index.gohtml", u)
	if err != nil {
		log.Fatal("Unable to execute template on index.gohtml")
		return
	}
}

func bar(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	err := tpl.ExecuteTemplate(w, "bar.gohtml", u)
	if err != nil {
		log.Fatal("Unable to parse bar.html")
		return
	}
}
