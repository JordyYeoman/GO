package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var pl = fmt.Println

func stripSpaces(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	return str
}

func main() {
	pl("Enter your favourite movie")
	reader := bufio.NewReader(os.Stdin) // Similar to C++, take user input
	movieName, err := reader.ReadString('\n')
	// Remove whitespace/newlines etc from string
	if err == nil {
		movieName = stripSpaces(movieName)
		pl("Hey ", movieName, "is my favourite movie also!")
	} else {
		log.Fatal(err)
	}
}
