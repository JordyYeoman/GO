package main

import (
	"encoding/json"
	"fmt"
	"github.com/JordyYeoman/GO/ToddMcleod_Go_Course_Udem/basic-mongo/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// curl -X POST -H "Content-Type: application/json" -d '{"Name":"James Bond","Gender":"male","Age":32,"Id":"777"}' http://127.0.0.1:8000/user

func main() {
	fmt.Println("Yo yo")
	r := mux.NewRouter()
	r.HandleFunc("/index", index).Methods("GET")
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/user", createUser).Methods("POST")
	r.HandleFunc("/user/{id}", getUser)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("unable to start server")
	}
}

func getUser(w http.ResponseWriter, res *http.Request) {
	vars := mux.Vars(res)
	u := models.User{
		Name:   "Tony Stark",
		Gender: "male",
		Age:    34,
		Id:     vars["id"],
	}

	//Marshal into json
	uj, _ := json.Marshal(u)

	// Write content, status code, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	_, err := fmt.Fprintf(w, "%s\n", uj)
	if err != nil {
		fmt.Println("Err writing json payload")
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// composite literal - type and curly braces
	u := models.User{}

	// encode/decode for sending/receiving JSON to/from a stream
	json.NewDecoder(r.Body).Decode(&u)

	// Change Id
	u.Id = "007"

	// marshal/unmarshal for having JSON assigned to a variable
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	_, err := fmt.Fprintf(w, "%s\n", uj)
	if err != nil {
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	// TODO: write code to delete user
	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprint(w, "Write code to delete user\n")
}

func index(w http.ResponseWriter, r *http.Request) {
	s := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Index</title>
</head>
<body>
<a href="/user/9872309847">GO TO: http://localhost:8080/user/9872309847</a>
</body>
</html>
	`
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}

func home(w http.ResponseWriter, res *http.Request) {
	_, err := fmt.Fprint(w, "Welcome homie!\n")
	if err != nil {
		log.Fatal("unable to load home route", err)
	}
}
