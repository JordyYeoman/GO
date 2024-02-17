package main

import (
	"fmt"
	"strconv"
	"unicode"
)

// Remember to use uppercase declaration if you want to export
var TeamNames = map[string]bool{
	"Richmond":               true,
	"Carlton":                true,
	"Geelong":                true,
	"Collingwood":            true,
	"North Melbourne":        true,
	"West Coast":             true,
	"Port Adelaide":          true,
	"Brisbane Lions":         true,
	"Melbourne":              true,
	"Western Bulldogs":       true,
	"Gold Coast":             true,
	"Sydney":                 true,
	"Greater Western Sydney": true,
	"Adelaide":               true,
	"Hawthorn":               true,
	"Essendon":               true,
	"St Kilda":               true,
	"Fremantle":              true,
}

func GetFinalScore(str string) int {
	var tempScore string
	for _, char := range str {
		if unicode.IsDigit(char) {
			tempScore += string(char)
		} else {
			break // Stop iteration if non-digit character encountered
		}
	}

	s, err := strconv.Atoi(tempScore)
	if err != nil {
		fmt.Println("Can't convert this to an int!")
	} else {
		fmt.Println(s)
	}

	return s
}
