package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Funcs(fm).ParseGlob("one.gohtml"))
}

// To inject functions into the template, we need to use a funcMap.
// "key" is the function name to call in the template
// "value" is the actual function we pass in.

// Get the first 3 letters of any string
func firstThree(s string) string {
	s = strings.TrimSpace(s)
	s = s[:3]
	return s
}
func fullSend(s string) string {
	fmt.Println(s)
	return s
}

var fm = template.FuncMap{
	"uc":       strings.ToUpper,
	"ft":       firstThree,
	"fullSend": fullSend,
}

type Person struct {
	Name string
	Age  int
}

type Car struct {
	Brand string
	Model string
	Year  int
}

type PageData struct {
	Cars   []Car
	People []Person
}

func main() {
	person1 := Person{
		Name: "Jordy",
		Age:  30,
	}

	person2 := Person{
		Name: "Amara",
		Age:  33,
	}

	person3 := Person{
		Name: "River",
		Age:  1,
	}

	car1 := Car{
		Brand: "Toyota",
		Model: "Landcruiser",
		Year:  1997,
	}

	car2 := Car{
		Brand: "Toyota",
		Model: "Rav4",
		Year:  2022,
	}

	people := []Person{person1, person2, person3}
	cars := []Car{car1, car2}

	pageData := PageData{
		cars,
		people,
	}
	err := tpl.ExecuteTemplate(os.Stdout, "one.gohtml", pageData)
	if err != nil {
		log.Fatal(err)
	}
}
