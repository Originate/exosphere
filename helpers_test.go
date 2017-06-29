package main

import "os/exec"

// Runs the given command, stores the output in "output"
func run(commands []string) (string, error) {
	command := exec.Command(commands[0], commands[1:]...)
	outputArray, err := command.CombinedOutput()
	output = string(outputArray)
	return output, err
}
