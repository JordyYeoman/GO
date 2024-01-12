package main

import "fmt"

func main() {
	// Do not assume these will be ordered -
	// GO makes no promises that these will be in any order.
	intMap := make(map[string]int)

	intMap["One"] = 1
	intMap["Two"] = 2
	intMap["Three"] = 3
	intMap["Four"] = 4
	intMap["Five"] = 5

	for key, value := range intMap {
		fmt.Println(key, value)
	}

	// Delete item from map
	// delete(intMap, "Four")

	// Check for a value in a map
	el, ok := intMap["Four"]
	if ok {
		fmt.Println(el, "is in map")
	} else {
		fmt.Println(el, "is NOT in map")
	}
}
