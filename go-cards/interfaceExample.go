package main

import "fmt"

type Bot interface {
	getGreeting() string
}
type EnglishBot struct{}
type SpanishBot struct{}

func main() {
	fmt.Println("Running Interface Example")

	eb := EnglishBot{}
	sb := SpanishBot{}

	PrintGreeting(eb)
	PrintGreeting(sb)
}

func PrintGreeting(b Bot) {
	fmt.Println(b.getGreeting())
}

func (EnglishBot) getGreeting() string {
	// very custom logic
	return "Hi there!"
}

func (SpanishBot) getGreeting() string {
	// very custom logic
	return "Hola!"
}
