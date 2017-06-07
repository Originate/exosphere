package main_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/DATA-DOG/godog"
	"github.com/Originate/exocom/go/exocom-mock"
	"github.com/Originate/exocom/go/structs"
)

type ServiceConfig struct {
	Type string `yaml:type`
}

func getServiceConfig() (*ServiceConfig, error) {
	config := &ServiceConfig{}
	configBytes, err := ioutil.ReadFile("service.yml")
	if err != nil {
		return nil, fmt.Errorf("Error reading service.yml", err)
	}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling service.yml", err)
	}
	return config, nil
}

func getRole() (string, error) {
	config, err := getServiceConfig()
	if err != nil {
		return "", err
	}
	return config.Type, nil
}

func newExocomMock(port int) *exocomMock.ExoComMock {
	exocom := exocomMock.New()
	go func() {
		err := exocom.Listen(port)
		if err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("Error listening on exocom", err))
		}
	}()
	return exocom
}

func FeatureContext(s *godog.Suite) {
	var exocom *exocomMock.ExoComMock
	var role string
	var serviceCommand *exec.Cmd
	var serviceCommandStdout, serviceCommandStderr io.ReadCloser
	port := 4100

	s.BeforeSuite(func() {
		var err error
		exocom = newExocomMock(port)
		role, err = getRole()
		if err != nil {
			panic(err)
		}
	})

	s.BeforeScenario(func(arg1 interface{}) {
		serviceCommand = nil
	})

	s.AfterScenario(func(interface{}, error) {
		exocom.Reset()
		if serviceCommand != nil {
			err := serviceCommand.Process.Kill()
			if err != nil {
				panic(fmt.Errorf("Error when killing the service command: %v", err))
			}
			stderr, err := ioutil.ReadAll(serviceCommandStderr)
			if err != nil {
				panic(fmt.Errorf("Error reading stderr for service command: %v", err))
			}
			if len(stderr) > 0 {
				panic(fmt.Errorf("Service command printed to stderr: %s", stderr))
			}
		}
	})

	s.AfterSuite(func() {
		err := exocom.Close()
		if err != nil {
			panic(fmt.Errorf("Error closing exocom", err))
		}
	})

	s.Step(`^an ExoCom server$`, func() error {
		return nil // Empty step as this is done in the BeforeSuite
	})

	s.Step(`^an instance of this service$`, func() error {
		serviceCommand = exec.Command("go", "run", "server.go")
		env := os.Environ()
		env = append(env, fmt.Sprintf("EXOCOM_PORT=%d", port), fmt.Sprintf("ROLE=%d", role))
		serviceCommand.Env = env
		var err error
		serviceCommandStdout, err = serviceCommand.StdoutPipe()
		if err != nil {
			return err
		}
		serviceCommandStderr, err = serviceCommand.StderrPipe()
		if err != nil {
			return err
		}
		err = serviceCommand.Start()
		if err != nil {
			return err
		}
		return exocom.WaitForConnection()
	})

	s.Step(`^receiving the "([^"]*)" command$`, func(name string) error {
		return exocom.Send(structs.Message{Name: name})
	})

	s.Step(`^this service replies with a "([^"]*)" message$`, func(name string) error {
		_, err := exocom.WaitForMessageWithName(name)
		return err
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
		StopOnFailure: true,
		Paths:         paths,
	})

	os.Exit(status)
}
