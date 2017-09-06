package testHelpers

import (
	"context"
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
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/docker/tools"
	execplus "github.com/Originate/go-execplus"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
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

func checkoutTemplate(cwd, templateName string) error {
	_, filePath, _, _ := runtime.Caller(0)
	src := path.Join(path.Dir(filePath), "..", "example-templates", templateName)
	dest := path.Join(cwd, "tmp", templateName)
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

func killTestContainers(dockerComposeDir, appDir string) error {
	dockerComposeProjectName := composebuilder.GetDockerComposeProjectName(appDir)
	mockLogger := application.NewLogger([]string{}, []string{}, ioutil.Discard)
	cleanProcess, err := compose.KillAllContainers(compose.BaseOptions{
		DockerComposeDir: dockerComposeDir,
		LogChannel:       mockLogger.GetLogChannel("feature-test"),
		Env:              []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", dockerComposeProjectName)},
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

func hasNetwork(dockerClient *client.Client, networkName string) (bool, error) {
	ctx := context.Background()
	networks, err := dockerClient.NetworkList(ctx, dockerTypes.NetworkListOptions{})
	if err != nil {
		return false, err
	}
	for _, network := range networks {
		if network.Name == networkName {
			return true, nil
		}
	}
	return false, nil
}

func listContainersInNetwork(dockerClient *client.Client, networkName string) ([]string, error) {
	containers := []string{}
	result, err := dockerClient.NetworkInspect(context.Background(), networkName, false)
	if err != nil {
		return []string{}, err
	}
	for _, container := range result.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil

}

func hasContainer(dockerClient *client.Client, containerName string) (bool, error) {
	containerNames, err := tools.ListRunningContainers(dockerClient)
	if err != nil {
		return false, err
	}
	for _, container := range containerNames {
		if container == containerName {
			return true, nil
		}
	}
	return false, nil
}

func waitForContainer(dockerClient *client.Client, containerName string) error {
	for {
		hasContainer, err := hasContainer(dockerClient, containerName)
		if err != nil {
			return err
		}
		if hasContainer {
			return nil
		}
	}
}

func runComposeInNetwork(command, network, path string) (*execplus.CmdPlus, error) {
	process := execplus.NewCmdPlus("docker-compose", "-p", network, command)
	process.SetDir(path)
	return process, process.Start()
}
