package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("Full Send")

	x := 7
	y := 7.11
	// Using %T, we can log out the type of variable passed into the PrintF statement.
	fmt.Printf("%T\n", x)
	fmt.Printf("%T\n", y)

	// Practice Questions
	// #1
	handsOnQuestions1()
	// #2
	handsOnQuestions2()
}

type Square struct {
	length float64
	width  float64
	units  string
}

type Circle struct {
	radius float64
	units  string
}

// Shape - Encompassing all shapes into an interface,
// they require the 'area()' method.
type Shape interface {
	area() float64
}

// Using a receiver to attach a method to the Shape interface.
func getArea(s Shape) float64 {
	return s.area()
}

func (c Circle) area() float64 {
	return math.Pi * (c.radius * c.radius)
}

// Here we are assigning the function to a type using a 'receiver'.
func (s Square) area() float64 {
	return s.width * s.length
}

func handsOnQuestions1() {
	testSquare := Square{
		11,
		10,
		"cm",
	}

	// Calling the function attached to the type
	fmt.Printf("Area of your square is %+v%+v\n", testSquare.area(), testSquare.units)

	testCircle := Circle{
		10,
		"mm",
	}

	// Calling the function attached to the type
	fmt.Printf("Area of your circle is %+v%+v\n", math.Round(testCircle.area()*100)/100, testCircle.units)

	// Now we have created the interface for all 'Shapes' and attached the method using a receiver, we can
	// simply call that method and pass in our 'Shape' to get the area.
	fmt.Printf("Area of circle is %+v\n", getArea(testCircle))
	fmt.Printf("Area of square is %+v\n", getArea(testSquare))
}

type Person struct {
	firstName string
	lastName  string
}

type SecretAgent struct {
	Person
}

func (p Person) pSpeak() {
	fmt.Printf("pSpeak\n")
}

func (s SecretAgent) saSpeak() {
	fmt.Printf("saSpeak\n")
}

func handsOnQuestions2() {
	p := Person{
		"Amara",
		"Yeoman",
	}
	secretAgent := SecretAgent{
		Person{
			"Jordy",
			"Yeoman",
		},
	}

	fmt.Println(p.firstName)
	fmt.Println(secretAgent.firstName)
	// Run speak method of Person type
	p.pSpeak()
	// Run speak method of SecretAgent type
	secretAgent.pSpeak()        // Works
	secretAgent.Person.pSpeak() // Also works - more verbose?
	secretAgent.saSpeak()
}
