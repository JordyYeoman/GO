package main

import "fmt"

func main() {
	myTotal := sumMany(2, 5, 6, 12, 515, 12, 5, 2)
	fmt.Println("Total: ", myTotal)
}

// When you have dynamic amount of properties to perform functions on
func sumMany(nums ...int) int {
	total := 0

	for _, x := range nums {
		total = total + x
	}

	return total
}
