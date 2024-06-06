package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

type Event struct {
	name     string
	location string
}

func EventRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/list", getAllEvents)
	return r
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	e := Event{"Burning Man", "Black Rock City"}
	fmt.Println(e)
	respondWithJSON(w, 200, e)
}
