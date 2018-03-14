package main

import (
	"bufio"
	"fmt"
	"os"

	"bytes"
	"github.com/atotto/clipboard"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	for text != "q" { // break the loop if text == "q"
		fmt.Print("Enter your text: ")
		scanner.Scan()
		text = scanner.Text()
		ltcText := ltcWalker(text)
		fmt.Println("Text with LTC: ", ltcText)
		clipboard.WriteAll(ltcText)
	}
}

func ltcWalker(input string) string {
	stringParts := strings.Split(input, " ")

	ltcString := bytes.NewBufferString("")
	for _, val := range stringParts {
		ltcString.WriteString(ltc(val) + " ")
	}

	return strings.Trim(ltcString.String(), " ")
}

func ltc(input string) string {
	if len(input) <= 2 {
		return input
	}

	return fmt.Sprintf("%c%d%c", input[0], len(input)-2, input[len(input)-1])
}
