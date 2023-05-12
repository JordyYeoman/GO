package main

import (
	"bufio"
	"fmt"
	"math/rand"
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
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Passing 8 to rand.Int will give a rand number between 0 and 8
	firstNumber := rand.Intn(8) + 2
	secondNumber := rand.Intn(8) + 2
	subtraction := rand.Intn(8) + 2
	var answer int

	// Take user input
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Guess the number game!")
	fmt.Println("================================================")
	fmt.Println("")
	fmt.Println("Guess a number between 1 and 10", userInputPrompt)
	reader.ReadString('\n')

	// Play game
	fmt.Println("Multiply your number by", firstNumber, userInputPrompt)
	reader.ReadString('\n')

	fmt.Println("Multiply the result by", secondNumber, userInputPrompt)
	reader.ReadString('\n')

	fmt.Println("Divide the result by your original number", userInputPrompt)
	reader.ReadString('\n')
	
	fmt.Println("Now subtract", subtraction, userInputPrompt)
	reader.ReadString('\n')

	// Return answer to user
	answer = (firstNumber * secondNumber) - subtraction
	println("Answer is...");
	println("Drum roll please...")
	runNTimes(2)
	println(answer) 
}