package main

import (
	"fmt"
	"mnc-fullstack-technical-test/tahap-1/shared"
)

func main() {
	reader := shared.NewInputReader()

	fmt.Print("Input total belanja: ")
	total, ok := reader.ReadMoney()
	if !ok {
		fmt.Println("invalid input, must be number")
		return
	}

	fmt.Print("Input jumlah bayar: ")
	paid, ok := reader.ReadMoney()
	if !ok {
		fmt.Println("invalid input, must be number")
		return
	}

	if paid < total {
		fmt.Println("False, kurang bayar")
		return
	}

	change := paid - total
	rounded := shared.RoundDownToHundred(change)

	fmt.Printf("\nKembalian yang harus diberikan kasir: %s\n", shared.FormatNumber(change))
	fmt.Printf("Dibulatkan menjadi %s\n", shared.FormatNumber(rounded))
	fmt.Println("Pecahan uang:")

	for i := 0; i < len(shared.Denominations); i++ {
		denom := shared.Denominations[i]
		if rounded >= denom {
			count := rounded / denom
			rounded = rounded % denom

			if denom >= 1000 {
				fmt.Printf("%d lembar %s\n", count, shared.FormatNumber(denom))
			} else {
				fmt.Printf("%d koin %s\n", count, shared.FormatNumber(denom))
			}
		}
	}
}
