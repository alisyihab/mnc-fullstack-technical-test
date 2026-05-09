package shared

import (
	"strconv"
	"strings"
)

var Denominations = []int{
	100000,
	50000,
	20000,
	10000,
	5000,
	2000,
	1000,
	500,
	200,
	100,
}

func RoundDownToHundred(amount int) int {
	return amount - (amount % 100)
}

func FormatNumber(n int) string {
	s := strconv.Itoa(n)

	if n < 1000 {
		return s
	}

	var result []string

	for len(s) > 3 {
		result = append([]string{s[len(s)-3:]}, result...)
		s = s[:len(s)-3]
	}

	result = append([]string{s}, result...)

	return strings.Join(result, ".")
}

func ParseFormattedNumber(value string) (int, error) {
	clean := strings.ReplaceAll(value, ".", "")
	return strconv.Atoi(clean)
}
