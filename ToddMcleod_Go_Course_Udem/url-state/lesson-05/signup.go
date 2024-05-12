package main

import (
	"github.com/google/uuid"
	"log"
	"net/http"
)

func signup(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	// process form submission
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")

		// username taken?
		if _, ok := dbUsers[un]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}

		// create session
		sID, err := uuid.NewUUID()
		if err != nil {
			log.Fatal("Unable to create new user id")
		}
		c := &http.Cookie{
			Name:  "session-id",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = un

		// store user in dbUsers
		u := User{un, p, f, l, 29}
		dbUsers[un] = u

		// redirect
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	err := tpl.ExecuteTemplate(w, "signup.gohtml", nil)
	if err != nil {
		log.Fatal("Unable to execute signup template")
		return
	}
}
