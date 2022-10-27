package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

// Hello returns a greeting for the named person.
func Hello(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	return message
}

// When two or more consecutive named fun params share a type, you can omit all but the last type
func addThreeNumbers(x, y, z int) int {
	return x + y + z
}

// A function can return any numbers of results // tuple??
func swap(x, y string) (string, string) {
	return y, x
}

// Naked return
// A return statement without arguments returns the named return values. This is known as a "naked" return.
func split(sum float32) (x, y float32) {
	x = sum * 4 / 9
	y = sum - x
	return
}

func declareMaster() {
	var i, j, k = 1, 2, 3
	fmt.Println(i, j, k)
}

func shortDecl() int {
	k := 21 // Short ass declaration boy-O
	return k
}

func typeConversion() {
	var i int = 10
	var f float64 = float64(i) + 0.123
	fmt.Println(f)
	var u uint = uint(f)
	fmt.Println(u)
	z := 20
	y := float64(z) + 0.123
	fmt.Println(y)
	x := uint(y)
	fmt.Println(x)
}

func observeTypeConversion() {
	v := 42 // change me!
	fmt.Printf("v is of type %T\n", v)
	y := 32.21 // change me!
	fmt.Printf("v is of type %T\n", y)
}

func babiesFirstForLoop() {
	for i := 0.00; i <= 0.05; i += 0.01 {
		fmt.Print(i)
		fmt.Println(" Loop City ")
	}
}

func babiesFirstIfCheck() {
	x := "GO"
	z := "NO"
	if z == x {
		fmt.Println("z is the same as x")
	} else {
		fmt.Println("z is not GO")
	}
}

// Like for, the if statement can start with a short statement to execute before the condition.
// Variables declared by the statement are only in scope until the end of the if.
func pow(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim { // :O how cool is this shit!
		return v
	}
	return lim
}

// Calc sqr root
func calcSqrRoot(x float64) float64 {
	z := 1.0

	for i := 0; i <= 10; i++ {
		// Loop until closest approx to sqr root is found
		z -= (z*z - x) / (2 * z)
	}

	return z
}

func checkRuntime() {
	fmt.Print("Go runs on ")
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Printf("%s.\n", os)
	}
}

func howLongUntilSaturday() {
	fmt.Println("When's Saturday?")
	today := time.Now().Weekday()
	switch time.Saturday {
	case today + 0:
		fmt.Println("Today.")
	case today + 1:
		fmt.Println("Tomorrow.")
	case today + 2:
		fmt.Println("In two days.")
	default:
		fmt.Println("Too far away.")
	}
}

// simple use case for nonconditional switch
// Switch without a condition is the same as switch true.
func checkTime() {
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning!")
	case t.Hour() < 17:
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}
}

func deferShit() {
	defer fmt.Println("deferred shit")

	fmt.Println("yeah nah, yeahhh nahhhhhh")
}

// Go has pointers. A pointer holds the memory address of a value.
// The type *T is a pointer to a T value. Its zero value is nil.
// The & operator generates a pointer to its operand.
func pointersBebe() {
	i, j := 42, 2701

	p := &i         // point to i - remeber that & generates the pointer to i.
	fmt.Println(*p) // read i through the pointer
	*p = 21         // set i through the pointer
	fmt.Println(i)  // see the new value of i

	p = &j         // point to j
	*p = *p / 37   // divide j through the pointer
	fmt.Println(j) // see the new value of j
}

// Struct fields can be accessed through a struct pointer.
type Vertex struct {
	X int
	Y float32
}

func structPointers() {
	v := Vertex{47, 20198124.12}
	p := &v      // generate a pointer to the v variable
	p.X = 1e9    // Assign a new value to the v struct value X, using the pointer
	(*p).Y = 2.1 // Reassign using c like syntax (not as clean)
	fmt.Println(v)
}

func main() {
	fmt.Println(Hello("Jordy"))
	fmt.Println(addThreeNumbers(10, 20, 30)) // 60
	fmt.Println(swap("Hey", "You"))          // "You" "Hey"
	fmt.Println(split(20))                   // 8, 12
	declareMaster()
	fmt.Println(shortDecl())
	typeConversion()
	observeTypeConversion()
	babiesFirstForLoop()
	babiesFirstIfCheck()
	fmt.Println(
		pow(3, 2, 10),
		pow(3, 3, 20),
	)
	fmt.Println(calcSqrRoot(10))
	checkRuntime()
	howLongUntilSaturday()
	checkTime()
	deferShit()
	pointersBebe()
	structPointers()
}

// TIL
// float32 ==  Single-precision floating-point format (sometimes called FP32 or float32) is a computer number format, usually occupying 32 bits in computer memory
