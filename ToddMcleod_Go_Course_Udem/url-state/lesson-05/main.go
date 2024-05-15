package main

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"text/template"
)

type User struct {
	UserName string
	Password []byte
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

	// Create dummy user
	bs, _ := bcrypt.GenerateFromPassword([]byte("shakenNotStirred"), bcrypt.MinCost)
	dbUsers["test@test.com"] = User{"test@test.com", bs, "John", "Doe", 37}
}

func main() {
	fmt.Println("User Unique Session ID + lookup")

	http.HandleFunc("/", index)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
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

func login(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	var u User
	// handle login
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		pw := req.FormValue("password")
		u, ok := dbUsers[un]
		if !ok {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}

		// Is the password correct
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(pw))
		if err != nil {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}

		// create session cookie
		sID, err := uuid.NewUUID()
		if err != nil {
			log.Fatal("Login failed")
			return
		}
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}

		http.SetCookie(w, c)
		dbSessions[c.Value] = un
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	err := tpl.ExecuteTemplate(w, "login.gohtml", u)
	if err != nil {
		log.Fatal("Login failed")
		return
	}
}
