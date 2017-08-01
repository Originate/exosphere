package websocketbridge_test

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
	"github.com/Originate/exocom/go/exorelay"
	"github.com/Originate/exocom/go/structs"
	"github.com/Originate/exocom/go/utils"
	"github.com/Originate/exocom/go/websocketbridge"
	"github.com/phayes/freeport"
	uuid "github.com/satori/go.uuid"

	"github.com/gorilla/websocket"
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
	var clientPort int
	var clientURL string
	var exocom *exocomMock.ExoComMock
	var websocketBridgeInstance *websocketbridge.WebsocketBridge
	var clientWebsocket *websocket.Conn
	var clientCount int
	var outgoingMessageId string
	var savedError error

	s.BeforeScenario(func(interface{}) {
		exocomPort = freeport.GetPort()
		clientPort = freeport.GetPort()
		clientURL = fmt.Sprintf("ws://localhost:%v", clientPort)
		outgoingMessageId = ""
		savedError = nil
		clientWebsocket = nil
		clientCount = 0
	})

	s.AfterScenario(func(interface{}, error) {
		err := websocketBridgeInstance.Close()
		if err != nil {
			panic(err)
		}
		err = exocom.Close()
		if err != nil {
			panic(err)
		}
	})

	s.Step(`a websocket bridge connected to Exocom`, func() error {
		exocom = newExocom(exocomPort)
		websocketBridgeInstance = &websocketbridge.WebsocketBridge{}
		config := exorelay.Config{Host: "localhost", Port: exocomPort, Role: "websocket-test"}
		return websocketBridgeInstance.Open(config, clientPort)
	})

	s.Step(`receiving this message from the client:`, func(payloadStr *gherkin.DocString) error {
		var err error
		clientWebsocket, err = utils.ConnectWithRetry(clientURL, 100)
		if err != nil {
			return err
		}
		var message structs.Message
		err = json.Unmarshal([]byte(payloadStr.Content), &message)
		if err != nil {
			return err
		}
		outgoingMessageId = message.ID
		marshaledMessage, err := json.Marshal(message)
		return clientWebsocket.WriteMessage(websocket.TextMessage, marshaledMessage)
	})

	s.Step(`^receiving a message "([^"]*)" from (\d+) different clients$`, func(messageName string, count int) error {
		clientCount = count
		var err error
		for i := 0; i < clientCount; i++ {
			clientWebsocket, err = utils.ConnectWithRetry(clientURL, 100)
			if err != nil {
				return err
			}
			message := structs.Message{Name: messageName, ID: uuid.NewV4().String()}
			marshaledMessage, err := json.Marshal(message)
			if err != nil {
				return err
			}
			err = clientWebsocket.WriteMessage(websocket.TextMessage, marshaledMessage)
			if err != nil {
				return err
			}
		}
		return nil
	})

	s.Step(`receiving this message from exocom:`, func(payloadStr *gherkin.DocString) error {
		var message structs.Message
		err := json.Unmarshal([]byte(payloadStr.Content), &message)
		if err != nil {
			return err
		}
		_, err = exocom.WaitForConnection()
		if err != nil {
			return err
		}
		if message.SessionID != "" {
			receivedMessages := exocom.GetReceivedMessages()
			message.SessionID = receivedMessages[len(receivedMessages)-1].SessionID
		}
		return exocom.Send(message)
	})

	s.Step(`^those clients have different sessionIds$`, func() error {
		err := utils.WaitFor(func() bool {
			return len(exocom.GetReceivedMessages()) == 3
		}, "Expected exocom to receive 3 messages")
		if err != nil {
			return err
		}
		receivedMessages := exocom.GetReceivedMessages()
		sessionIdSet := make(map[string]bool)
		for _, message := range receivedMessages {
			if message.Name != "exocom.register-service" {
				sessionIdSet[message.SessionID] = true
			}
		}
		if len(sessionIdSet) != clientCount {
			return fmt.Errorf("Expected request to have %d clients but got %d", clientCount, len(sessionIdSet))
		}
		return nil
	})

	s.Step(`websocket bridge makes the websocket request:`, func(payloadStr *gherkin.DocString) error {
		var expectedMessage structs.Message
		err := json.Unmarshal([]byte(payloadStr.Content), &expectedMessage)
		if err != nil {
			return err
		}
		actualMessage, err := exocom.WaitForMessageWithName(expectedMessage.Name)
		if err != nil {
			return err
		}
		//ID and SessionId get generated on the fly, so we can't know them beforehand
		expectedMessage.ID = actualMessage.ID
		expectedMessage.SessionID = actualMessage.SessionID
		if !reflect.DeepEqual(actualMessage, expectedMessage) {
			return fmt.Errorf("Expected request to equal %s but got %s", expectedMessage, actualMessage)
		}
		return nil
	})

	s.Step(`websocket bridge sends this to the client:`, func(payloadStr *gherkin.DocString) error {
		_, bytes, err := clientWebsocket.ReadMessage()
		if err != nil {
			return err
		}
		var actualMessage structs.Message
		err = json.Unmarshal(bytes, &actualMessage)
		if err != nil {
			return err
		}
		var expectedMessage structs.Message
		err = json.Unmarshal([]byte(payloadStr.Content), &expectedMessage)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(actualMessage, expectedMessage) {
			return fmt.Errorf("Expected request to equal %s but got %s", expectedMessage, actualMessage)
		}
		return nil
	})

	s.Step(`^the websocket bridge does not send a message to the client$`, func() error {
		var bytes []byte
		clientWebsocket.SetReadDeadline(time.Now().Add(time.Second * time.Duration(1)))
		_, bytes, err := clientWebsocket.ReadMessage()
		if err == nil {
			return fmt.Errorf("Expected timeout error")
		}
		if strings.Contains(err.Error(), "i/o timeout") {
			if bytes != nil {
				var actualMessage structs.Message
				err = json.Unmarshal(bytes, &actualMessage)
				if err != nil {
					return err
				}
				return fmt.Errorf("Expected no message, but received %v", actualMessage)
			}
			return nil
		} else {
			return err
		}
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
