package userInput

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func Ask(reader *bufio.Reader, query string) string {
	fmt.Print(query)
	answer, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to reader user input: %s", err)
	}
	answer = strings.TrimSpace(answer)
	if len(answer) == 0 {
		fmt.Println("(expect a non-empty string)\n")
		return Ask(reader, query)
	}
	return answer
}

func Choose(reader *bufio.Reader, query string, options []string) string {
	fmt.Println(query)
	if len(options) == 0 {
		log.Fatal(fmt.Errorf("no options found\n"))
	}
	for i, option := range options {
		fmt.Printf("%v. %v\n", i+1, option)
	}
	fmt.Print(fmt.Sprintf("[expect a number between 1 and %v]: ", len(options)))
	answer, err := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	if err != nil {
		panic(err)
	}
	chosenNumber, err := strconv.Atoi(answer)
	if err != nil || !(0 <= chosenNumber-1 && chosenNumber-1 < len(options)) {
		fmt.Printf("error: expected a number between 1 and %v]\n\n", len(options))
		return Choose(reader, query, options)
	}
	return options[chosenNumber-1]
}
