package main

import "os/exec"

// Runs the given command, stores the output in "output"
func run(commands []string) (err error) {
	command := exec.Command(commands[0], commands[1:]...)
	command.Dir = testRoot
	outputArray, err := command.CombinedOutput()
	output = string(outputArray)
	return
}
