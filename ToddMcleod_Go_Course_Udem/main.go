package main

import "fmt"

func main() {
	fmt.Println("Full Send")

	x := 7
	y := 7.11
	// Using %T, we can log out the type of variable passed into the PrintF statement.
	fmt.Printf("%T\n", x)
	fmt.Printf("%T\n", y)

	// Practice Questions #1
	handsOnQuestions1()
}

type Square struct {
	length float64
	width  float64
	units  string
}

// Here we are assigning the function to a type using a 'receiver'.
func (s Square) getAreaForShape() float64 {
	return s.width * s.length
}

type Circle struct {
}

func handsOnQuestions1() {
	testSquare := Square{
		11,
		10,
		"cm",
	}

	// Calling the function attached to the type
	fmt.Printf("Area of your square is %+v%+v\n", testSquare.getAreaForShape(), testSquare.units)
}
