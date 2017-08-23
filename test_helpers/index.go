package testHelpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/docker/compose"
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
	src := path.Join(path.Dir(filePath), "..", "example-apps", appName)
	dest := path.Join(cwd, "tmp", appName)
	err := os.RemoveAll(dest)
	if err != nil {
		return err
	}
	return CopyDir(src, dest)
}

func createEmptyApp(appName, cwd string) (string, error) {
	parentDir := os.TempDir()
	cmdPlus := execplus.NewCmdPlus("exo", "create")
	cmdPlus.SetDir(parentDir)
	if err := cmdPlus.Start(); err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Failed to create '%s' application", appDir))
	}
	fields := []string{"AppName", "AppDescription", "AppVersion", "ExocomVersion"}
	inputs := []string{appName, "Empty test application", "1.0.0", "0.24.0"}
	for i, field := range fields {
		if err := cmdPlus.WaitForText(field, time.Second*5); err != nil {
			return "", err
		}
		if _, err := cmdPlus.StdinPipe.Write([]byte(inputs[i] + "\n")); err != nil {
			return "", err
		}
	}
	return path.Join(parentDir, appName), nil
}

func killTestContainers(dockerComposeDir string) error {
	mockLogger := application.NewLogger([]string{}, []string{}, ioutil.Discard)
	cleanProcess, err := compose.KillAllContainers(compose.BaseOptions{
		DockerComposeDir: dockerComposeDir,
		LogChannel:       mockLogger.GetLogChannel("feature-test"),
	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Output:%s", cleanProcess.GetOutput()))
	}
	err = cleanProcess.Wait()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Output:%s", cleanProcess.GetOutput()))
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
	return cmdPlus.WaitForText("all services online", time.Minute*2)
}

func enterInput(row *gherkin.TableRow) error {
	field, input := row.Cells[0].Value, row.Cells[1].Value
	if err := childCmdPlus.WaitForText(field, time.Second*5); err != nil {
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
