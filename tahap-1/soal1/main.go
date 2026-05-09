package main

import (
	"fmt"
	"mnc-fullstack-technical-test/tahap-1/shared"
)

func main() {
	reader := shared.NewInputReader()

	fmt.Print("Input total string: ")

	n, ok := reader.ReadInt()
	if !ok {
		fmt.Println("invalid input, must be integer")
		return
	}

	if n <= 0 {
		fmt.Println("false")
		return
	}

	norms := make([]string, n)

	for i := 0; i < n; i++ {
		fmt.Printf("Input string %d: ", i+1)

		line, ok := reader.ReadLine()
		if !ok {
			fmt.Printf("invalid input, must be %d strings\n", n)
			return
		}

		norms[i] = shared.NormalizeString(line)
	}

	type groupInfo struct {
		count     int
		firstSeen int
	}

	groups := make(map[string]*groupInfo)

	for i := 0; i < n; i++ {
		key := norms[i]

		if info, exists := groups[key]; exists {
			info.count++
		} else {
			groups[key] = &groupInfo{
				count:     1,
				firstSeen: i,
			}
		}
	}

	bestKey := ""
	bestCount := 1
	bestFirstSeen := n

	for key, info := range groups {
		if info.count > bestCount {
			bestKey = key
			bestCount = info.count
			bestFirstSeen = info.firstSeen
			continue
		}

		if info.count == bestCount && info.count > 1 && info.firstSeen < bestFirstSeen {
			bestKey = key
			bestFirstSeen = info.firstSeen
		}
	}

	if bestCount <= 1 {
		fmt.Println("false")
		return
	}

	first := true
	for i := 0; i < n; i++ {
		if norms[i] == bestKey {
			if !first {
				fmt.Print(" ")
			}
			fmt.Print(i + 1)
			first = false
		}
	}

	fmt.Println()
}