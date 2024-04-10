package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// Important!!
// type handler interface {
//	ServeHttp(ResponseWriter, *Request)
// }

type hotdog int

// Extending the above hotdog type to include
func (m hotdog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal("Could not parse form: ", err)
	}

	fmt.Println(r.Form)

	err = tpl.ExecuteTemplate(w, "index.gohtml", r.Form)
	if err != nil {
		log.Fatal("Unable to execute template.")
	}
	fmt.Fprintf(w, "Any code you want in this cheeky bugger")
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.gohtml"))
}

func main() {
	var d hotdog

	err := http.ListenAndServe(":8080", d)
	if err != nil {
		log.Fatal("Server Died")
	}

	return
}
