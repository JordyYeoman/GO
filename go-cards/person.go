package main

import "fmt"

type Person struct {
	firstName string
	lastName  string
	age       int
}

func DoStuff() {
	j := Person{
		firstName: "Jordy",
		lastName:  "Yeoman",
		age:       29,
	}

	jPointer := &j
	fmt.Printf("Jordy: %+v\n", j)
	jPointer.updateName("Maggie")
	j.logPerson()
}

func (p *Person) updateName(newName string) {
	(*p).firstName = newName
}

func (p *Person) logPerson() {
	fmt.Printf("Maggie: %+v\n", p)
}
