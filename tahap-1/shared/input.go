package shared

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type InputReader struct {
	scanner *bufio.Scanner
}

func NewInputReader() *InputReader {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024), 1024*1024)

	return &InputReader{
		scanner: scanner,
	}
}

func (r *InputReader) ReadLine() (string, bool) {
	if !r.scanner.Scan() {
		return "", false
	}
	return strings.TrimSpace(r.scanner.Text()), true
}

func (r *InputReader) ReadInt() (int, bool) {
	line, ok := r.ReadLine()
	if !ok {
		return 0, false
	}

	value, err := strconv.Atoi(line)
	if err != nil {
		return 0, false
	}

	return value, true
}

func (r *InputReader) ReadMoney() (int, bool) {
	line, ok := r.ReadLine()
	if !ok {
		return 0, false
	}

	value, err := ParseFormattedNumber(line)
	if err != nil {
		return 0, false
	}

	return value, true
}
