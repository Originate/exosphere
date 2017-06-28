package userInput

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func Ask(reader *bufio.Reader, query string, required bool) string {
	fmt.Print(query)
	answer, err := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	if err != nil {
		panic(err)
	}
	if len(answer) == 0 && required {
		fmt.Println("expect a non-empty string")
		return Ask(reader, query, required)
	}
	return answer
}

func Choose(reader *bufio.Reader, query string, options []string) string {
	fmt.Print(query)
	if len(options) == 0 {
		panic(fmt.Errorf("no options found\n"))
	}
	for i, option := range options {
		fmt.Printf("%v. %v\n", i+1, option)
	}
	answer, err := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	if err != nil {
		panic(err)
	}
	chosenNumber, err := strconv.Atoi(answer)
	if err != nil || !(0 <= chosenNumber-1 && chosenNumber-1 < len(options)) {
		fmt.Printf("expect a number between 1 and %v\n", len(options))
		return Choose(reader, query, options)
	}
	return options[chosenNumber-1]
}
