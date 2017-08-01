package exoservice_test

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/Originate/exocom/go/exocom-mock"
	"github.com/Originate/exocom/go/exorelay"
	"github.com/Originate/exocom/go/exoservice"
	"github.com/Originate/exocom/go/exoservice/test-fixtures"
	"github.com/Originate/exocom/go/structs"
	"github.com/phayes/freeport"
)

func newExocom(port int) *exocomMock.ExoComMock {
	exocom := exocomMock.New()
	go func() {
		err := exocom.Listen(port)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	return exocom
}

func FeatureContext(s *godog.Suite) {
	var exocomPort int
	var exocom *exocomMock.ExoComMock
	var exoService *exoservice.ExoService
	var testFixture exoserviceTestFixtures.TestFixture

	s.BeforeScenario(func(interface{}) {
		exocomPort = freeport.GetPort()
	})

	s.AfterScenario(func(interface{}, error) {
		err := exoService.Close()
		if err != nil {
			panic(err)
		}
		err = exocom.Close()
		if err != nil {
			panic(err)
		}
	})

	s.Step(`^I connect the "([^"]*)" test fixture$`, func(name string) error {
		exocom = newExocom(exocomPort)
		testFixture = exoserviceTestFixtures.Get(name)
		config := exorelay.Config{
			Host: "localhost",
			Port: exocomPort,
			Role: "test-service",
		}
		exoService = &exoservice.ExoService{}
		err := exoService.Connect(config)
		if err != nil {
			return err
		}
		go func() {
			err := exoService.ListenForMessages(testFixture.GetMessageHandler())
			if err != nil {
				panic(err)
			}
		}()
		return nil
	})

	s.Step(`^receiving a "([^"]*)" message(?: with (?:id "([^"]*)")?(?:(?: and )?sessionId "([^"]*)")?)?$`, func(name, id, sessionId string) error {
		message := structs.Message{
			ID:        id,
			Name:      name,
			SessionID: sessionId,
		}
		_, err := exocom.WaitForConnection()
		if err != nil {
			return err
		}
		return exocom.Send(message)
	})

	s.Step(`^it sends a "([^"]*)" message(?: as a reply to the message with (?:id "([^"]*)")?(?:(?: and )?sessionId "([^"]*)")?)?$`, func(name, id, sessionId string) error {
		actualMessage, err := exocom.WaitForMessageWithName(name)
		if err != nil {
			return err
		}
		if actualMessage.ResponseTo != id {
			return fmt.Errorf("Expected message to be a response to %s but got %s", id, actualMessage.ResponseTo)
		}
		if actualMessage.SessionID != sessionId {
			return fmt.Errorf("Expected message to be a response to have sessionId %s but got %s", sessionId, actualMessage.SessionID)
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
