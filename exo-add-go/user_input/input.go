package userInput

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func Ask(reader *bufio.Reader, query string) string {
	fmt.Print(query)
	answer, err := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	if len(answer) == 0 {
		fmt.Println("expect a non-empty string")
		return Ask(reader, query, required)
	}
	return answer
}
