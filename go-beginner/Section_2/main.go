package main

import "time"

func runNTimes(n int) {
	for i := 0; i < n; i++ {
		println("....")
		time.Sleep(1 * time.Second)
	}
}

func main() {
	println("Say WHAT AGAIN MOTHERFKER")
	println("I DARE YOU, I DOUBLE DARE YOU!")
	runNTimes(3)
	println("What?!")
}