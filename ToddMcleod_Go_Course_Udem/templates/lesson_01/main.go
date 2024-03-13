package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
	tpl, err := template.ParseFiles("tpl.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	nf, err := os.Create("index.html")
	if err != nil {
		log.Println("Error creating new file", err)
	}
	defer func(nf *os.File) {
		err := nf.Close()
		if err != nil {
			log.Fatal("Unable to close file")
		}
	}(nf)

	err = tpl.Execute(os.Stdout, nil)
	if err != nil {
		log.Fatal(err)
	}
}
