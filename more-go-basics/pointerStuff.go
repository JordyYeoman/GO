package main

import "fmt"

func zero(x int) {
	x = 0
}

func pointerZero(varPtr *int) {
	*varPtr = 0
}

// pointer is represented using a *
//

func pointerStuff() {
	x := 5
	y := 10
	zero(x)
	pointerZero(&y) // value at that memory address
	fmt.Println(y)
	fmt.Println(x) // x is still 5
}
