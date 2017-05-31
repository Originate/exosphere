package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/DATA-DOG/godog"
	"github.com/Originate/exocom/go/exocom-mock"
	"github.com/Originate/exocom/go/structs"
	"github.com/Originate/exocom/go/utils"
)

type ServiceConfig struct {
	Type string `yaml:type`
}

func getRole() string {
	configBytes, err := ioutil.ReadFile("service.yml")
	if err != nil {
		panic(fmt.Errorf("Error reading service.yml", err))
	}
	config := ServiceConfig{}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		panic(fmt.Errorf("Error unmarshaling service.yml", err))
	}
	return config.Type
}

func newExocom(port int) *exocomMock.ExoComMock {
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
	port := 4100

	s.BeforeSuite(func() {
		exocom = newExocom(port)
		role = getRole()
	})

	s.BeforeScenario(func(arg1 interface{}) {
		serviceCommand = nil
	})

	s.AfterScenario(func(interface{}, error) {
		exocom.Reset()
		if serviceCommand != nil {
			err := serviceCommand.Process.Kill()
			if err != nil {
				panic(fmt.Errorf("Error when killing the service command", err))
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
		return serviceCommand.Start()
	})

	s.Step(`^receiving the "([^"]*)" command$`, func(name string) error {
		message := structs.Message{Name: name}
		err := utils.WaitFor(func() bool { return exocom.HasConnection() }, "nothing connected to exocom")
		if err != nil {
			return err
		}
		return exocom.Send(message)
	})

	s.Step(`^this service replies with a "([^"]*)" message$`, func(name string) error {
		err := exocom.WaitForReceivedMessagesCount(2)
		if err != nil {
			return err
		}
		actualMessage := exocom.ReceivedMessages[1]
		if actualMessage.Name != name {
			return fmt.Errorf("Expected message to have name %s but got %s", name, actualMessage.Name)
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
		StopOnFailure: true,
		Paths:         paths,
	})

	os.Exit(status)
}
