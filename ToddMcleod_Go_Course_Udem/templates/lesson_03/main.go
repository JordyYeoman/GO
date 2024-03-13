package main

import (
	"log"
	"os"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("one.gohtml"))
}

func main() {
	err := tpl.ExecuteTemplate(os.Stdout, "one.gohtml", 42)
	if err != nil {
		log.Fatal(err)
	}
}
