package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	fmt.Println("Encoding script running")
	s := "Obadiah Stane: [shouting] Tony Stark was able to build this in a cave! With a box of scraps!"
	s64 := base64.StdEncoding.EncodeToString([]byte(s))

	fmt.Println(len(s))
	fmt.Println(len(s64))
	fmt.Println(s)
	fmt.Println(s64)
}
