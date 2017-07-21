package appSetup

import (
	"bufio"
	"os"
	"strings"

	"github.com/Originate/exosphere/exo-go/src/os_helpers"
)

func getCommands(dockerfilePath string) ([]string, error) {
	var commands []string
	if osHelpers.FileExists(dockerfilePath) {
		f, err := os.Open(dockerfilePath)
		if err != nil {
			return commands, err
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "RUN") {
				commands = append(commands, strings.Split(line, "RUN ")[1])
			}
		}
		if err := f.Close(); err != nil {
			return commands, err
		}
	}
	return commands, nil
}
