package testHelpers

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/exo-go/src/application"
	"github.com/Originate/exosphere/exo-go/src/docker"
	execplus "github.com/Originate/go-execplus"
	"github.com/pkg/errors"
)

const validateTextContainsErrorTemplate = `
Expected:

%s

to include

%s
	`

// CheckoutApp copies the example app appName to cwd
func CheckoutApp(cwd, appName string) error {
	_, filePath, _, _ := runtime.Caller(0)
	src := path.Join(path.Dir(filePath), "..", "..", "example-apps", appName)
	dest := path.Join(cwd, "tmp", appName)
	err := os.RemoveAll(dest)
	if err != nil {
		return err
	}
	return CopyDir(src, dest)
}

func killTestContainers(dockerComposeDir string) error {
	_, pipeWriter := io.Pipe()
	mockLogger := application.NewLogger([]string{}, []string{}, pipeWriter)
	cleanProcess, err := docker.KillAllContainers(dockerComposeDir, mockLogger.GetLogChannel("feature-test"))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Output:%s", cleanProcess.Output))
	}
	err = cleanProcess.Wait()
	if err != nil {
		fmt.Printf("Error:%s\nOutput:%s\n", err, cleanProcess.Output)
	}
	return nil
}

func runApp(cwd, appName string) error {
	appDir = path.Join(cwd, "tmp", appName)
	cmdPlus := execplus.NewCmdPlus("exo", "run") // nolint gas
	cmdPlus.SetDir(appDir)
	err := cmdPlus.Start()
	if err != nil {
		return err
	}
	return cmdPlus.WaitForText("all services online", time.Minute*5)
}

func enterInput(row *gherkin.TableRow) error {
	field, input := row.Cells[0].Value, row.Cells[1].Value
	if err := childCmdPlus.WaitForText(field, time.Second); err != nil {
		return err
	}
	_, err := childCmdPlus.StdinPipe.Write([]byte(input + "\n"))
	return err
}

func validateTextContains(haystack, needle string) error {
	if strings.Contains(haystack, needle) {
		return nil
	}
	return fmt.Errorf(validateTextContainsErrorTemplate, haystack, needle)
}
