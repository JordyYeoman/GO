package main

import "log"

func IsEvenOrOdd(listOfNumbers []int) {
	for _, num := range listOfNumbers {
		if num%2 == 0 {
			log.Printf("%v, Is Even\n", num)
		} else {
			log.Printf("%v, Is Odd\n", num)
		}
	}
}
