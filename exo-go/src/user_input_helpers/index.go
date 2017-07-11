package userInputHelpers

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Ask asks the user to enter an answer for the query and
// returns the trimmed answer
func Ask(reader *bufio.Reader, query string) string {
	fmt.Print(query)
	answer, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to reader user input: %s", err)
	}
	answer = strings.TrimSpace(answer)
	if len(answer) == 0 {
		fmt.Print("(expect a non-empty string)\n\n")
		return Ask(reader, query)
	}
	return answer
}

// Choose asks the user to select an option from the given
// list of options, and returns the selected option
func Choose(reader *bufio.Reader, query string, options []string) string {
	fmt.Println(query)
	if len(options) == 0 {
		log.Fatal(fmt.Errorf("no options found"))
	}
	for i, option := range options {
		fmt.Printf("%v - %v\n", i+1, option)
	}
	fmt.Print(fmt.Sprintf("[expect a number between 1 and %v]: ", len(options)))
	answer, err := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	if err != nil {
		panic(err)
	}
	chosenNumber, err := strconv.Atoi(answer)
	if err != nil || !(0 <= chosenNumber-1 && chosenNumber-1 < len(options)) {
		fmt.Printf("error: expected a number between 1 and %v\n\n", len(options))
		return Choose(reader, query, options)
	}
	return options[chosenNumber-1]
}

// Confirm asks the user to answer "yes" or "no" to the given query
func Confirm(reader *bufio.Reader, query string) bool {
	fmt.Printf("%s (y or n): ", query)
	answer, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read user input: %s", err)
	}
	answer = strings.TrimSpace(answer)
	switch answer {
	case "y":
		return true
	case "n":
		return false
	default:
		fmt.Println(`error: expected "y" or "n"`)
		return Confirm(reader, query)
	}
}
