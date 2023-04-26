package main

import (
	"fmt"
	"myapp/doctor"
	"myapp/jarvis"
)

func main() {
	//
	sendItVar := "sendItVar"
	fmt.Println("SEND ITTTT")
	getSome(sendItVar)
	//
	// Testing external package
	fmt.Println(doctor.Intro())
	// doMethod := doctor.Intro()
	// Testing external package 2
	fmt.Println(jarvis.Intro())
}

func getSome(input string) {
	fmt.Println(input)
}