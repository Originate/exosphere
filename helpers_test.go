package main

import (
	"fmt"
	shutil "github.com/termie/go-shutil"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func commitAllChanges() {
	err := run([]string{"git", "add", "-A"})
	checkText(err, output)
	err = run([]string{"git", "commit", "-m", "changes"})
	checkText(err, output)
}

// returns whether the given command list contains a color argument
func containsColorArgument(commands []string) bool {
	return strings.Contains(strings.Join(commands, " "), "--color=")
}

func createTempDir() (path string) {
	dir, err := ioutil.TempDir("", "morula")
	check(err)
	return dir
}

// Creates a subproject of the given type,
// with the given name,
// in the given test workspace
func createTestProject(template string, name string) {
	check(shutil.CopyTree(
		filepath.Join("features", "examples", template),
		filepath.Join(testRoot, name),
		nil))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkText(e error, text string) {
	if e != nil {
		fmt.Println(text)
		panic(e)
	}
}

func makeCrossPlatformCommand(command string) string {
	if os.PathSeparator == '\\' {
		command = strings.Replace(command, "/", "\\\\", -1)
		command = strings.Replace(command, "bin\\spec", "bin\\spec.cmd", -1)
	}
	return command
}

// Runs the given command, stores the output in "output"
func run(commands []string) (err error) {
	if commands[0] == "morula" && !containsColorArgument(commands) {
		commands = append(commands, "--color=false")
	}
	command := exec.Command(commands[0], commands[1:]...)
	command.Dir = testRoot
	outputArray, err := command.CombinedOutput()
	output = string(outputArray)
	return
}

// Returns the project names in the given human-friendly project name list
func splitProjectNames(projectNames string) (result []string) {
	for _, projectName := range strings.Split(projectNames, "and") {
		result = append(result, strings.Trim(projectName, " \""))
	}
	return
}

func switchBranch(branchName string) {
	err := run([]string{"git", "checkout", "-b", branchName})
	checkText(err, output)
}
