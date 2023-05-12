package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func runNTimes(n int) {
	for i := 0; i < n; i++ {
		println("....")
		time.Sleep(1 * time.Second)
	}
}

const userInputPrompt = "and press enter when ready"

func main() {
	println("Say WHAT AGAIN MOTHERFKER")
	println("I DARE YOU, I DOUBLE DARE YOU!")
	// runNTimes(3)
	println("What?!")

	firstNumber := 2
	secondNumber := 7
	subtraction := 3
	var answer int

	// Take user input
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Guess the number game!")
	fmt.Println("================================================")
	fmt.Println("")
	fmt.Println("Guess a number between 1 and 10", userInputPrompt)
	reader.ReadString('\n')

	// Play game
	fmt.Println("Multiply your number by ", firstNumber, userInputPrompt)
	reader.ReadString('\n')

	fmt.Println("Multiply the result by ", secondNumber, userInputPrompt)
	reader.ReadString('\n')

	fmt.Println("Divide the result by your original number ", userInputPrompt)
	reader.ReadString('\n')
	
	fmt.Println("Now subtract  ", subtraction, userInputPrompt)
	reader.ReadString('\n')

	// Return answer to user
}