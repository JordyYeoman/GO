package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	filePath := os.Args[1]

	fmt.Println("OS arguments: ")
	fmt.Println("File path: ", filePath)

	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		fmt.Println("Error opening file")
	}

	r := bufio.NewReader(f)

	// Section 2
	for {
		line, _, err := r.ReadLine()
		if len(line) > 0 {
			fmt.Printf("ReadLine: %q\n", line)
		}
		if err != nil {
			break
		}
	}
}
