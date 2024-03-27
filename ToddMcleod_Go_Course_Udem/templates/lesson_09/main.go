package main

import (
	"log"
	"os"
	"text/template"
)

type Player struct {
	Name         string
	Age          int
	AvgDisposals float64
}

type Team struct {
	Players []Player
	Name    string
}

var tpl *template.Template

func doubleIt(n float64) float64 {
	return n * 2
}

func squareNum(n float64) float64 {
	return n * n
}

var fm = template.FuncMap{
	"double": doubleIt,
	"sqr":    squareNum,
}

func init() {
	tpl = template.Must(template.New("").Funcs(fm).ParseGlob("tpl.gohtml"))
}

func main() {
	p1 := Player{
		"Jordy",
		29,
		27.7,
	}
	p2 := Player{
		"Sean",
		29,
		14.2,
	}
	p3 := Player{
		"Luke",
		29,
		29.2,
	}
	team := Team{
		Players: []Player{
			p1, p2, p3,
		},
		Name: "Parkerville",
	}

	err := tpl.ExecuteTemplate(os.Stdout, "tpl.gohtml", team)
	if err != nil {
		log.Fatal(err)
	}
}
