package main

import "fmt"

type Triangle struct {
	base   float64
	height float64
}
type Square struct {
	SideLength float64
}
type Shape interface {
	getArea() float64
}

func PrintArea(shapeName string, s Shape) {
	fmt.Println(shapeName+" area:", s.getArea())
}

func (t Triangle) getArea() float64 {
	return 0.5 * t.height * t.base
}

func (s Square) getArea() float64 {
	return s.SideLength * s.SideLength
}

func main() {
	fmt.Println("Starting program")

	triangle := Triangle{base: 64.2, height: 27.1}
	square := Square{SideLength: 12.241}

	PrintArea("Triangle", triangle)
	PrintArea("Square", square)
}
