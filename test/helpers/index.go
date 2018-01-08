package helpers

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
	"github.com/Originate/exosphere/src/docker/tools"
	exoContext "github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/util"
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

const validateTextDoesNotContainErrorTemplate = `
Expected:

%s

not to include

%s
`

// GetTestApplicationDir returns the path to the test application with the given name
func GetTestApplicationDir(appName string) string {
	_, filePath, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filePath), "..", "fixtures", "applications", appName)
}

// CheckoutApp copies the example app into the given appDir
func CheckoutApp(appDir, appName string) error {
	src := GetTestApplicationDir(appName)
	err := os.RemoveAll(appDir)
	if err != nil {
		return err
	}
	return CopyDir(src, appDir)
}

// JoinStringSlices joins the given slices
func JoinStringSlices(slices ...[]string) []string {
	result := []string{}
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

func checkoutTemplate(templateDir, templateName string) error {
	_, filePath, _, _ := runtime.Caller(0)
	src := path.Join(path.Dir(filePath), "..", "fixtures", "service_templates", templateName)
	err := os.RemoveAll(templateDir)
	if err != nil {
		return err
	}
	return CopyDir(src, templateDir)
}

func createEmptyApp(appName string) (string, error) {
	parentDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}
	cmdPlus := execplus.NewCmdPlus("exo", "init")
	cmdPlus.SetDir(parentDir)
	if err := cmdPlus.Start(); err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Failed to create '%s' application", appDir))
	}
	fields := []string{"AppName", "ExocomVersion"}
	inputs := []string{appName, "0.24.0"}
	for i, field := range fields {
		if err := cmdPlus.WaitForText(field, time.Second*5); err != nil {
			return "", err
		}
		if _, err := cmdPlus.StdinPipe.Write([]byte(inputs[i] + "\n")); err != nil {
			return "", err
		}
	}
	return parentDir, nil
}

func killAppContainers(appDir string) error {
	writer := ioutil.Discard
	appContext, err := exoContext.GetAppContext(appDir)
	if err != nil {
		return err
	}
	return application.CleanContainers(appContext, writer)
}

func cleanApp(appDir string) error {
	doesExist, err := util.DoesFileExist(path.Join(appDir, "application.yml"))
	if err != nil {
		return err
	}
	if doesExist {
		cmdPlus := execplus.NewCmdPlus("exo", "clean")
		cmdPlus.SetDir(appDir)
		err := cmdPlus.Run()
		if err != nil {
			return errors.Wrapf(err, "'exo clean' errored with output:\n%s", cmdPlus.GetOutput())
		}
	}
	return nil
}

func runApp(appDir, textToWaitFor string) error {
	cmdPlus := execplus.NewCmdPlus("exo", "run")
	cmdPlus.SetDir(appDir)
	err := cmdPlus.Start()
	if err != nil {
		return err
	}
	return cmdPlus.WaitForText(textToWaitFor, time.Minute*2)
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

func runComposeInNetwork(command, network, path, filename string, env []string) (*execplus.CmdPlus, error) {
	process := execplus.NewCmdPlus("docker-compose", "-p", network, "-f", filename, command)
	process.SetDir(path)
	process.AppendEnv(env)
	return process, process.Start()
}

func isEmptyFile(filePath string) (bool, error) {
	stat, err := os.Stat(filePath)
	if err == nil {
		return stat.Size() == 0, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
