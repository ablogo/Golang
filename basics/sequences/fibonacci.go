package sequences

import "fmt"

var memo = make(map[int]int)

func Fibonacci(number int) {
	a, b := 0, 1

	for number >= 0 {
		fmt.Printf("%d ", a)
		a, b = b, a+b
		number -= 1
	}
}

func Recursive(number int, a int, b int) int {
	number--
	if number == 0 {
		return b
	} else {
		a, b = b, a+b
		fmt.Printf("%d ", Recursive(number, a, b))
		return a
	}
}

func RecursiveMathApproach(n int) int {
	if n <= 1 {
		return n
	}
	return RecursiveMathApproach(n-1) + RecursiveMathApproach(n-2)
}

func RecursiveMathMemoization(n int) int {
	v, ok := memo[n]
	if ok {
		return v
	}

	if n <= 1 {
		memo[n] = n
		return n
	}
	memo[n] = RecursiveMathMemoization(n-1) + RecursiveMathMemoization(n-2)
	return memo[n]
}
