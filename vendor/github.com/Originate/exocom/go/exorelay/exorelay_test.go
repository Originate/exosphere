package exorelay_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exocom/go/exocom-mock"
	"github.com/Originate/exocom/go/exorelay"
	"github.com/Originate/exocom/go/exorelay/test-fixtures"
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

// Cucumber step definitions
// nolint gocyclo
func FeatureContext(s *godog.Suite) {
	var exocomPort int
	var exocom *exocomMock.ExoComMock
	var exoInstance *exorelay.ExoRelay
	var outgoingMessageId string
	var savedError error
	var testFixture exorelayTestFixtures.TestFixture

	s.BeforeScenario(func(interface{}) {
		exocomPort = freeport.GetPort()
		outgoingMessageId = ""
		savedError = nil
		testFixture = nil
	})

	s.AfterScenario(func(interface{}, error) {
		err := exoInstance.Close()
		if err != nil {
			panic(err)
		}
		err = exocom.Close()
		if err != nil {
			panic(err)
		}
	})

	s.Step(`^an ExoRelay with the role "([^"]*)"$`, func(role string) error {
		exoInstance = &exorelay.ExoRelay{
			Config: exorelay.Config{
				Host: "localhost",
				Port: exocomPort,
				Role: role,
			},
		}
		return nil
	})

	s.Step(`^(an ExoRelay instance that is connected to Exocom|ExoRelay connects to Exocom)$`, func() error {
		exocom = newExocom(exocomPort)
		return exoInstance.Connect()
	})

	s.Step(`^it registers by sending the message "([^"]*)" with payload:$`, func(expectedName string, payloadStr *gherkin.DocString) error {
		message, err := exocom.WaitForMessageWithName(expectedName)
		if err != nil {
			return err
		}
		var expectedPayload structs.MessagePayload
		err = json.Unmarshal([]byte(payloadStr.Content), &expectedPayload)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(message.Payload, expectedPayload) {
			return fmt.Errorf("Expected message payload to equal %s but got %s", expectedPayload, message.Payload)
		}
		return nil
	})

	s.Step(`^sending the message "([^"]*)"$`, func(name string) error {
		var err error
		outgoingMessageId, err = exoInstance.Send(exorelay.MessageOptions{Name: name})
		return err
	})

	s.Step(`^sending the message "([^"]*)" with the payload:$`, func(name string, payloadStr *gherkin.DocString) error {
		var payload structs.MessagePayload
		err := json.Unmarshal([]byte(payloadStr.Content), &payload)
		if err != nil {
			return err
		}
		outgoingMessageId, err = exoInstance.Send(exorelay.MessageOptions{Name: name, Payload: payload})
		return err
	})

	s.Step(`^sending the message "([^"]*)" with sessionId "([^"]*)"$`, func(name, sessionId string) error {
		var err error
		outgoingMessageId, err = exoInstance.Send(exorelay.MessageOptions{Name: name, SessionID: sessionId})
		return err
	})

	s.Step(`^trying to send an empty message$`, func() error {
		outgoingMessageId, savedError = exoInstance.Send(exorelay.MessageOptions{Name: ""})
		if savedError == nil {
			return fmt.Errorf("Expected ExoRelay to error but it did not")
		} else {
			return nil
		}
	})

	s.Step(`^ExoRelay makes the WebSocket request:$`, func(messageStr *gherkin.DocString) error {
		t := template.New("request")
		t, err := t.Parse(messageStr.Content)
		if err != nil {
			return err
		}
		var expectedMessageBuffer bytes.Buffer
		err = t.Execute(&expectedMessageBuffer, map[string]interface{}{"outgoingMessageId": outgoingMessageId})
		if err != nil {
			return err
		}
		var expectedMessage structs.Message
		err = json.Unmarshal(expectedMessageBuffer.Bytes(), &expectedMessage)
		if err != nil {
			return err
		}
		actualMessage, err := exocom.WaitForMessageWithName(expectedMessage.Name)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(actualMessage, expectedMessage) {
			return fmt.Errorf("Expected request to equal %s but got %s", expectedMessage, actualMessage)
		}
		return nil
	})

	s.Step(`^ExoRelay errors with "([^"]*)"$`, func(expectedErrorMessage string) error {
		actualErrorMessage := savedError.Error()
		if actualErrorMessage != expectedErrorMessage {
			return fmt.Errorf("Expected error to equal %s but got %s", expectedErrorMessage, actualErrorMessage)
		}
		return nil
	})

	s.Step(`^I setup the "([^"]*)" test fixture$`, func(name string) error {
		testFixture = exorelayTestFixtures.Get(name)
		testFixture.Setup(exoInstance)
		return nil
	})

	s.Step(`^receiving this message:$`, func(messageStr *gherkin.DocString) error {
		var message structs.Message
		err := json.Unmarshal([]byte(messageStr.Content), &message)
		if err != nil {
			return err
		}
		_, err = exocom.WaitForConnection()
		if err != nil {
			return err
		}
		return exocom.Send(message)
	})

	s.Step(`^the fixture receives a message with the name "([^"]*)" and the payload nil$`, func(messageName string) error {
		message, err := testFixture.WaitForMessageWithName(messageName)
		if err != nil {
			return err
		}
		actualPayload := message.Payload
		if actualPayload != nil {
			return fmt.Errorf("Expected payload to nil but got %s", actualPayload)
		}
		return nil
	})

	s.Step(`^the fixture receives a message with the name "([^"]*)" and the payload:$`, func(messageName string, payloadStr *gherkin.DocString) error {
		var expectedPayload structs.MessagePayload
		err := json.Unmarshal([]byte(payloadStr.Content), &expectedPayload)
		if err != nil {
			return err
		}
		message, err := testFixture.WaitForMessageWithName(messageName)
		if err != nil {
			return err
		}
		actualPayload := message.Payload
		if !reflect.DeepEqual(actualPayload, expectedPayload) {
			return fmt.Errorf("Expected payload to %s but got %s", expectedPayload, actualPayload)
		}
		return nil
	})

	s.Step(`^the fixture receives a message with the name "([^"]*)" and sessionId "([^"]*)"$`, func(messageName, sessionId string) error {
		message, err := testFixture.WaitForMessageWithName(messageName)
		if err != nil {
			return err
		}
		if message.SessionID != sessionId {
			return fmt.Errorf("Expected session id %s but got %s", sessionId, message.SessionID)
		}
		return nil
	})

	s.Step(`^Exocom is offline$`, func() error {
		return nil // noop
	})

	s.Step(`^ExoRelay and Exocom boot up simultaneously$`, func() error {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			exoInstance.Connect()
			wg.Done()
		}()
		time.Sleep(time.Duration(200) * time.Millisecond)
		exocom = newExocom(exocomPort)
		wg.Wait()
		return nil
	})

	s.Step(`^ExoRelay should (?:re)?connect to Exocom$`, func() error {
		_, err := exocom.WaitForConnection()
		if err != nil {
			return err
		}
		return nil
	})

	s.Step(`^Exocom goes offline momentarily`, func() error {
		return exocom.CloseConnection()
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
