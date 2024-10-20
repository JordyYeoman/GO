package main

import "fmt"

func main() {
	colors := map[string]string{
		"red":   "#ff0000",
		"pink":  "#FFC0CB",
		"green": "#00ff00",
	}

	printMap(colors)
	fmt.Println("Colors: ", colors)
}

func printMap(c map[string]string) {
	for colour, hex := range c {
		fmt.Println("Colour: ", colour, "Hex: ", hex)
	}
}
