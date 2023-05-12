package main

import (
	"fmt"
	"sort"
)

// basic types
var myInt int
var myInt16 int16
var myFloat float32

// aggregate types (array, struct)
var myTestArray [4]string

// interface types
type Car struct {
	NumberOfTyres int
	Luxury bool
	BucketSeats bool
	Make string
	Model string
	Year int
}

// Unsigned integer, which can only be assigned to positive numbers or 0
var myUInt uint

func basicTypes() {
	println("Understanding types")
	println("=============================================")

	// Basic types
	myInt = 12
	myInt16 = 12344	
	myFloat = 12.11

	println(myInt, myInt16, myFloat)

	// aggregated types
	// Property inside square brackets is the length of the array
	// When the array is created, it will be autofilled with all 0's
	var myStrings [3]string
	var myInts [4]int

	myStrings[0] = "Dog"
	myStrings[1] = "Cat"
	myStrings[2] = "Fish"

	myInts[0] = 1
	myInts[1] = 2
	myInts[2] = 3

	fmt.Println("First element in myStrings array of strings is:", myStrings[2])
	fmt.Println("First element in myInts array of strings is:", myInts[3])

	// Structs
	var myCar Car
	myCar.BucketSeats = true
	myCar.Luxury = true
	myCar.Make = "Tesla"
	myCar.Model = "Y"
	myCar.Year = 2021

	fmt.Printf("My Car is a %d %s %s", myCar.Year, myCar.Make, myCar.Model)
}

func main() {
	basicTypes()
	pointers()
	runSlice()
}

func runSlice() {
	var animals []string

	animals = append(animals, "dog")
	animals = append(animals, "cat")
	animals = append(animals, "fish")
	animals = append(animals, "snail")
	animals = append(animals, "horse")

	fmt.Println("animals:", animals)

	// loop over animals with for loop
	// Where '_' is the index & range is the length of array/slice
	loopElementsOfSlice(animals)

	// Get first element
	fmt.Println("Element 0 is:", animals[0])
	// Get first 2 elements
	fmt.Println("First 2 elements:", animals[0:2])

	// Check length of slice
	fmt.Println("Slice length is:", len(animals), "elements long")

	// Sort slice
	fmt.Println("Is it sorted??", sort.StringsAreSorted(animals))
	
	sort.Strings(animals)

	fmt.Println("Is animals slice sorted?", sort.StringsAreSorted(animals))
	fmt.Println("Sort animals:", animals)

	// Remove element from slice
	fmt.Println("Removing the 3rd element, which is:", animals[2])
	// TODO: Understand why the final element in the array is not actually deleted
	deleteFromSlice(animals, 2)
	fmt.Println("New Animals slice:", animals)
	loopElementsOfSlice(animals)
}

func deleteFromSlice(a []string, i int) []string {
	lengthOfSlice := len(a) - 1

	a[i] = a[lengthOfSlice]
	a[lengthOfSlice] = ""
	a = a[:lengthOfSlice]
	sort.Strings(a)
	return a
}

func loopElementsOfSlice(sliceEl []string) {
	for i, x := range sliceEl {
		fmt.Println(i, x)
	}
}

func pointers() {
	// Pointer is simply something that points to a location in memory
	var myInt int
	myInt = 10

	fmt.Println(myInt)

	// Example
	x := 10

	myFirstPointer := &x

	fmt.Println("x is a ", x)
	fmt.Println("myFirst Pointer is", myFirstPointer)

	// Go to the address in memory of where the value myFirstPointer is stored
	// and update the value
	*myFirstPointer = 24

	fmt.Println("x is now", x)

	changeValueOfPointer(*&myFirstPointer)

	fmt.Println("x is now..", x)
}

func changeValueOfPointer(num *int) {
	*num = 25
}