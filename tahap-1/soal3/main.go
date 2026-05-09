package main

import (
	"fmt"
	"mnc-fullstack-technical-test/tahap-1/shared"
)

func main() {
	reader := shared.NewInputReader()

	fmt.Print("Input string bracket: ")

	input, ok := reader.ReadLine()
	if !ok {
		fmt.Println("false")
		return
	}

	if len(input) < 1 || len(input) > 4096 {
		fmt.Println("false")
		return
	}

	stack := shared.NewRuneStack()

	pairs := map[rune]rune{
		'>': '<',
		'}': '{',
		']': '[',
	}

	openings := map[rune]bool{
		'<': true,
		'{': true,
		'[': true,
	}

	for _, char := range input {

		// opening bracket
		if openings[char] {
			stack.Push(char)
			continue
		}

		// closing bracket
		expectedOpening, exists := pairs[char]
		if !exists {
			fmt.Println("false")
			return
		}

		lastOpening, ok := stack.Pop()
		if !ok {
			fmt.Println("false")
			return
		}

		if lastOpening != expectedOpening {
			fmt.Println("false")
			return
		}
	}

	if !stack.IsEmpty() {
		fmt.Println("false")
		return
	}

	fmt.Println("true")
}
