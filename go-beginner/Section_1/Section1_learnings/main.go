package main

import (
	"bufio"
	"fmt"
	"myapp/doctor"
	"myapp/jarvis"
	"os"
	"strings"
	"time"
)

func main() {
	testDoctor := doctor.Intro()
	testJarvis := jarvis.Intro()

	fmt.Println(testDoctor)
	fmt.Println(testJarvis)

	// Reader
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Print("--> ")
		// Take in user input and assign to variable
		userInput, _ := reader.ReadString('\n')

		userInput = strings.Replace(userInput, "\r\n", "", -1)
		userInput = strings.Replace(userInput, "\n", "", -1)

		if userInput == "quit" {
			break
		} else {
			fmt.Println(doctor.Response(userInput))
		}
	}

	// Testing switch statements
	switchWithoutExpression()
	switchWithExpression("test")

	fmt.Println("Program terminated successfully.")
}

func switchWithoutExpression() {
	// A switch without a condition is the same as switch true.
	t := time.Now()
    switch {
    case t.Hour() < 12:
        fmt.Println("It's before noon")
    default:
        fmt.Println("It's after noon")
    }
}

func switchWithExpression(input string) {
	switch input {
		case "yes":
			fmt.Println("Yes match")
		case "no":
			fmt.Println("No match")
		case "test":
			fmt.Println("test match")
		default:
			fmt.Println("No cases match")
	}
}