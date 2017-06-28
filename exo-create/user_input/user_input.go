package userInput

import (
	"bufio"
	"fmt"
	"strings"
)

func Ask(reader *bufio.Reader, query string, defaultVal string, required bool) string {
	if len(defaultVal) > 0 {
		query = query + "(" + defaultVal + ") "
	}
	fmt.Print(query)
	answer, err := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	if err != nil {
		panic(err)
	}
	if len(answer) == 0 {
		if len(defaultVal) > 0 {
			answer = defaultVal
		} else if required {
			fmt.Println("expect a non-empty string")
			return Ask(reader, query, defaultVal, required)
		}
	}
	return answer
}
