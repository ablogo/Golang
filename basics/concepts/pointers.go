package concepts

import "fmt"

// A pointer is a variable that stores the memory address of another variable.
// To create a pointer use the "&" operator to get the address of the variable.
// "&" address of
// Use the "*" operator to declare a pointer type or access the value of that address (pointer).
// "*" value at (dereferencing)

func Pointers() {
	var pointerType *int
	x := 15
	p := &x
	pointerType = &x
	fmt.Printf("variable value (x): %d \n", x)
	fmt.Printf("Pointer (p) address: %d \n", p)
	fmt.Printf("Pointer value (p): %d \n", *p)

	// change the value of that address
	*p = 20
	fmt.Printf("New value. x = %d, p = %d \n", x, *p)
	fmt.Printf("Pointer type: %d \n", *pointerType)
}

type Counter struct {
	Count int
}

// Value receiver
func (c Counter) GetCount() int {
	return c.Count
}

// Pointer receiver
func (c *Counter) Increment() {
	c.Count++
}

// Not modify the caller
func (c Counter) Decrement() {
	c.Count--
}

// Modify the caller
func (c *Counter) Reset() {
	c.Count = 0
}

func StructExample() {
	counter := Counter{Count: 0}
	counter.Increment()
	counter.Decrement()
	counter.Increment()
	counter.Decrement()
	fmt.Printf("Counter should be 0. Counter value: %d \n", counter.GetCount())

	counter.Reset()
	fmt.Printf("Reset counter, current value: %d \n", counter.GetCount())

}
