package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("one.gohtml"))
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
	// Slice of strings
	//family := []string{"Jordy", "Amara", "River"}

	// Map
	//family := map[string]string{
	//	"Kalgoorlie": "Jordy",
	//	"Armadale":   "Amara",
	//	"Midland":    "River",
	//}

	// struct
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
	//err := tpl.ExecuteTemplate(os.Stdout, "one.gohtml", family)
	err := tpl.ExecuteTemplate(os.Stdout, "one.gohtml", pageData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
