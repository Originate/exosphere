package exocomMock

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Originate/exocom/go/structs"
	"github.com/Originate/exocom/go/utils"
	"golang.org/x/net/websocket"
)

// ExoComMock is a mock implementation of ExoRelay,
// to be used for testing
type ExoComMock struct {
	ReceivedMessages []structs.Message
	server           http.Server
	socket           *websocket.Conn
}

// New creates a new ExoComMock instance
func New() *ExoComMock {
	result := new(ExoComMock)
	result.server = http.Server{
		Handler: websocket.Handler(result.websocketHandler),
	}
	return result
}

// Close takes this ExoComMock instance offline
func (e *ExoComMock) Close() error {
	return e.server.Close()
}

// Listen brings this ExoComMock instance online
func (e *ExoComMock) Listen(port int) error {
	e.server.Addr = fmt.Sprintf(":%d", port)
	return e.server.ListenAndServe()
}

// HasConnection returns whether or not a socket is connected
func (e *ExoComMock) HasConnection() bool {
	return e.socket != nil
}

// Reset closes and nils the socket and clears all received messages
func (e *ExoComMock) Reset() {
	if e.HasConnection() {
		err := e.socket.Close()
		if err != nil {
			log.Fatal(err)
		}
		e.socket = nil
	}
	e.ReceivedMessages = []structs.Message{}
}

// Send sends the given message to the connected socket
func (e *ExoComMock) Send(message structs.Message) error {
	if !e.HasConnection() {
		return fmt.Errorf("Nothing connected to exocom")
	}
	serializedBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return websocket.Message.Send(e.socket, serializedBytes)
}

// WaitForConnection waits for a socket to connect
func (e *ExoComMock) WaitForConnection() error {
	return utils.WaitFor(func() bool {
		return e.HasConnection()
	}, "Expected a socket to connect to exocom")
}

// WaitForReceivedMessagesCount waits the received messages count to equal the given count
func (e *ExoComMock) WaitForReceivedMessagesCount(count int) error {
	return utils.WaitFor(func() bool {
		return len(e.ReceivedMessages) >= count
	}, fmt.Sprintf("Expected exocom to recieve %d messages but only has %d:\n%v", count, len(e.ReceivedMessages), e.ReceivedMessages))
}

// Helpers

func (e *ExoComMock) readMessage(socket *websocket.Conn) (structs.Message, error) {
	var bytes []byte
	if err := websocket.Message.Receive(socket, &bytes); err != nil {
		return structs.Message{}, err
	}

	var unmarshaled structs.Message
	err := json.Unmarshal(bytes, &unmarshaled)
	if err != nil {
		return structs.Message{}, err
	}

	return unmarshaled, nil
}

func (e *ExoComMock) websocketHandler(socket *websocket.Conn) {
	e.socket = socket
	for {
		message, err := e.readMessage(socket)
		if socket != e.socket {
			break // socket was closed via the Reset function
		} else if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading message from websocket:", err)
			break
		} else {
			e.ReceivedMessages = append(e.ReceivedMessages, message)
		}
	}
}
