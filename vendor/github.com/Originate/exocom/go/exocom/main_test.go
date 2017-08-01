package main_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exocom/go/exocom/test_helpers"
	"github.com/phayes/freeport"
	"github.com/pkg/errors"
)

// nolint gocyclo
func FeatureContext(s *godog.Suite) {
	var cmd *exec.Cmd
	var cmdOutput string
	var cmdStdout io.ReadCloser
	var cmdStderr io.ReadCloser
	var shouldCmdBeKilled bool
	var server *http.Server
	var exocomPort int
	var services map[string]*testHelpers.MockService

	createService := func(serviceName string) error {
		service := testHelpers.NewMockService(exocomPort, serviceName)
		err := service.Connect()
		if err != nil {
			return err
		}
		err = testHelpers.WaitForText(cmdStdout, fmt.Sprintf("'%s' registered", serviceName))
		if err != nil {
			return err
		}
		services[serviceName] = service
		return nil
	}

	startExocom := func(env []string) error {
		var err error
		cmd = exec.Command("exocom")
		cmd.Env = env
		cmdStdout, err = cmd.StdoutPipe()
		if err != nil {
			return err
		}
		cmdStderr, err = cmd.StderrPipe()
		if err != nil {
			return err
		}
		return cmd.Start()
	}

	s.BeforeScenario(func(arg1 interface{}) {
		cmd = nil
		cmdOutput = ""
		cmdStdout = nil
		cmdStderr = nil
		shouldCmdBeKilled = true
		server = nil
		exocomPort = freeport.GetPort()
		services = map[string]*testHelpers.MockService{}
	})

	s.AfterScenario(func(arg1 interface{}, arg2 error) {
		if cmd != nil && cmd.Process != nil && shouldCmdBeKilled {
			err := cmd.Process.Kill()
			if err != nil {
				panic(err)
			}
		}
		if server != nil {
			err := server.Close()
			if err != nil {
				panic(err)
			}
		}
		for _, service := range services {
			err := service.Close()
			if err != nil {
				panic(err)
			}
		}
	})

	s.Step(`^starting ExoCom$`, func() error {
		return startExocom([]string{
			fmt.Sprintf("SERVICE_ROUTES=%s", "[]"),
		})
	})

	s.Step(`^another service already uses port (\d+)$`, func(port int) error {
		server = testHelpers.StartMockServer(port)
		return nil
	})

	s.Step(`^I see "([^"]*)"$`, func(text string) error {
		return testHelpers.WaitForText(cmdStdout, text)
	})

	s.Step(`^it aborts with the message "([^"]*)"$`, func(text string) error {
		if !strings.Contains(cmdOutput, text) {
			return fmt.Errorf("Expected '%s' to contain '%s'", cmdOutput, text)
		}
		return nil
	})

	s.Step(`^starting ExoCom at port (\d+)$`, func(port int) error {
		exocomPort = port
		return startExocom([]string{
			fmt.Sprintf("PORT=%d", exocomPort),
			fmt.Sprintf("SERVICE_ROUTES=%s", "[]"),
		})
	})

	s.Step(`^trying to start ExoCom at port (\d+)$`, func(port int) error {
		var err error
		cmd = exec.Command("exocom")
		cmd.Env = []string{
			fmt.Sprintf("PORT=%d", port),
			fmt.Sprintf("SERVICE_ROUTES=%s", "[]"),
		}
		output, err := cmd.CombinedOutput()
		if err == nil {
			return fmt.Errorf("Expected exocom to fail but it didn't")
		}
		cmdOutput = string(output)
		shouldCmdBeKilled = false
		return nil
	})

	s.Step(`^an ExoCom instance configured with the routes:$`, func(docString *gherkin.DocString) error {
		return startExocom([]string{
			fmt.Sprintf("PORT=%d", exocomPort),
			fmt.Sprintf("SERVICE_ROUTES=%s", docString.Content),
		})
	})

	s.Step(`^a "([^"]*)" service connects and registers itself$`, func(serviceName string) error {
		return createService(serviceName)
	})

	s.Step(`^a running "([^"]*)" instance$`, func(serviceName string) error {
		return createService(serviceName)
	})

	s.Step(`^the "([^"]*)" service goes offline$`, func(serviceName string) error {
		return services[serviceName].Close()
	})

	s.Step(`^ExoCom should have the config:$`, func(servicesDocString *gherkin.DocString) error {
		var expected map[string]interface{}
		err := json.Unmarshal([]byte(servicesDocString.Content), &expected)
		if err != nil {
			return err
		}
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/config.json", exocomPort))
		if err != nil {
			return err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var actual map[string]interface{}
		err = json.Unmarshal(body, &actual)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Error unmarshalling: %s", string(body)))
		}
		if !reflect.DeepEqual(expected, actual) {
			return fmt.Errorf("Expected to equal %s but got %s", expected, actual)
		}
		return nil
	})
}

func TestMain(m *testing.M) {
	var paths []string
	var format string
	if len(os.Args) == 3 && os.Args[1] == "--" {
		format = "pretty"
		paths = append(paths, os.Args[2])
	} else {
		format = "progress"
		paths = append(paths, "features")
	}
	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:        format,
		NoColors:      false,
		StopOnFailure: false,
		Paths:         paths,
	})

	os.Exit(status)
}
