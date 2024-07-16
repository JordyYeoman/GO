package main

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func getUser(w http.ResponseWriter, req *http.Request) User {
	// get cookie
	c, err := req.Cookie("session-id")
	if err != nil {
		sID, _ := uuid.NewUUID()
		c = &http.Cookie{
			Name:  "session-id",
			Value: sID.String(),
		}
	}
	c.MaxAge = sessionLength
	http.SetCookie(w, c)

	var u User
	// if the user exists already, get user
	if session, ok := dbSessions[c.Value]; ok {
		session.lastActivity = time.Now()
		dbSessions[c.Value] = session
		u = dbUsers[session.un]
	}
	showSessions() // for testing
	return u
}

func alreadyLoggedIn(req *http.Request) bool {
	c, err := req.Cookie("session-id")
	if err != nil {
		return false
	}

	session := dbSessions[c.Value]
	_, ok := dbUsers[session.un]
	return ok
}

func cleanSessions() {
	fmt.Println("BEFORE CLEAN")
	showSessions()

	for k, v := range dbSessions {
		if time.Now().Sub(v.lastActivity) > (time.Second * 30) {
			delete(dbSessions, k)
		}
	}
	dbSessionsCleaned = time.Now()
	fmt.Println("AFTER CLEAN")
	showSessions()
}

func showSessions() {
	fmt.Println("=====================")
	for k, v := range dbSessions {
		fmt.Println(k, v.un)
	}
	fmt.Println(" ")
}
