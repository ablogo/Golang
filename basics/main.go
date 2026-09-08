package main

import (
	"basics/sequences"
	"fmt"
)

func main() {
	sequences.Fibonacci(8)
	fmt.Println("")
	sequences.Recursive(8, 0, 1)
	fmt.Println("")
	r := sequences.RecursiveMathApproach(8)
	fmt.Println(r)
	r = sequences.RecursiveMathMemoization(8)
	fmt.Println(r)

	fmt.Println("Hello ")
}
