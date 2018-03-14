package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/atotto/clipboard"
)

func main() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			s := <-signalChannel
			switch s {
			case syscall.SIGINT, syscall.SIGTERM:
				os.Exit(0)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	var text string
	for {
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
