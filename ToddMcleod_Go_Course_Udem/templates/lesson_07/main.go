package main

import (
	"log"
	"os"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("one.gohtml"))
}

type GameState struct {
	Score1  int
	Score2  int
	player1 string
	player2 string
}

type PageData struct {
	Game  GameState
	Count []string
}

func main() {
	xs := []string{"zero", "one", "two", "three"}
	newGame := GameState{
		24,
		15,
		"Jordy",
		"Amara",
	}
	pageData := PageData{
		newGame,
		xs,
	}

	err := tpl.ExecuteTemplate(os.Stdout, "one.gohtml", pageData)
	if err != nil {
		log.Fatal(err)
	}
}
