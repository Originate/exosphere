package exocomMock_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exocom/go/exocom-mock"
	"github.com/Originate/exocom/go/structs"
	"github.com/Originate/exocom/go/utils"
	"github.com/gorilla/websocket"
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

// nolint gocyclo
func FeatureContext(s *godog.Suite) {
	var exocomPort int
	var exocom *exocomMock.ExoComMock
	var socket *websocket.Conn

	s.BeforeScenario(func(interface{}) {
		exocomPort = freeport.GetPort()
	})

	s.AfterScenario(func(interface{}, error) {
		err := exocom.Close()
		if err != nil {
			panic(err)
		}
	})

	s.Step(`^an Exocom server$`, func() error {
		exocom = newExocom(exocomPort)
		return nil
	})

	s.Step(`^I have mocked it to reply to "([^"]*)" with "([^"]*)" and the payload:$`, func(requestName, replyName string, replyPayloadStr *gherkin.DocString) error {
		var replyPayload structs.MessagePayload
		err := json.Unmarshal([]byte(replyPayloadStr.Content), &replyPayload)
		if err != nil {
			return err
		}
		exocom.AddMockReply(requestName, exocomMock.MockReplyData{
			Name:    replyName,
			Payload: replyPayload,
		})
		return nil
	})

	s.Step(`^receiving a "([^"]*)" message with id "([^"]*)"$`, func(name, id string) error {
		var err error
		utils.Retry(5, func() bool {
			socket, _, err = websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%d", "localhost", exocomPort), nil)
			return err == nil
		})
		if err != nil {
			return err
		}
		serializedBytes, err := json.Marshal(structs.Message{
			Name: name,
			ID:   id,
		})
		if err != nil {
			return err
		}
		return socket.WriteMessage(websocket.TextMessage, serializedBytes)
	})

	s.Step(`^it sends a "([^"]*)" message as a reply "([^"]*)" with the payload$`, func(name, id string, payloadStr *gherkin.DocString) error {
		var expectedPayload structs.MessagePayload
		err := json.Unmarshal([]byte(payloadStr.Content), &expectedPayload)
		if err != nil {
			return err
		}
		deadline := time.Now().Add(time.Second)
		err = socket.SetReadDeadline(deadline)
		if err != nil {
			return err
		}
		_, bytes, err := socket.ReadMessage()
		if err != nil {
			return err
		}
		message, err := utils.MessageFromJSON(bytes)
		if err != nil {
			return err
		}
		if message.Name != name {
			return fmt.Errorf("Expected message name to equal '%s' but got '%s'", name, message.Name)
		}
		if message.ResponseTo != id {
			return fmt.Errorf("Expected message responseTo to equal '%s' but got '%s'", id, message.ResponseTo)
		}
		if !reflect.DeepEqual(message.Payload, expectedPayload) {
			return fmt.Errorf("Expected message payload to equal %s but got %s", expectedPayload, message.Payload)
		}
		return nil
	})

	s.Step(`^it does not send a message$`, func() error {
		deadline := time.Now().Add(time.Second)
		err := socket.SetReadDeadline(deadline)
		if err != nil {
			return err
		}
		_, bytes, err := socket.ReadMessage()
		if err == nil {
			return fmt.Errorf("Expected no message to be received but received: %s", string(bytes))
		}
		if strings.Contains(err.Error(), "i/o timeout") {
			return nil
		}
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
		StopOnFailure: false,
		Paths:         paths,
	})

	os.Exit(status)
}
