package exorelay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/Originate/exocom/go/structs"
	uuid "github.com/satori/go.uuid"

	"golang.org/x/net/websocket"
)

// Config contains the configuration values for ExoRelay instances
type Config struct {
	Host string
	Port int
	Role string
}

// MessageOptions contains the configuration values for ExoRelay instances
type MessageOptions struct {
	Name       string
	Payload    structs.MessagePayload
	ResponseTo string
}

// ExoRelay is the low level Go API to talk to Exocom
type ExoRelay struct {
	Config         Config
	socket         *websocket.Conn
	messageChannel chan structs.Message
}

// Connect brings an ExoRelay instance online
func (e *ExoRelay) Connect() error {
	socket, err := websocket.Dial(fmt.Sprintf("ws://%s:%d", e.Config.Host, e.Config.Port), "", "origin:")
	if err != nil {
		return err
	}
	e.socket = socket
	e.messageChannel = make(chan structs.Message)
	go e.listenForMessages()
	_, err = e.Send(MessageOptions{
		Name:    "exocom.register-service",
		Payload: map[string]interface{}{"clientName": e.Config.Role},
	})
	return err
}

// GetMessageChannel returns a channel which can be used read incoming messages
func (e *ExoRelay) GetMessageChannel() chan structs.Message {
	return e.messageChannel
}

// Send sends a message with the given options
func (e *ExoRelay) Send(options MessageOptions) (string, error) {
	id := uuid.NewV4().String()
	if options.Name == "" {
		return "", errors.New("ExoRelay#Send cannot send empty messages")
	}
	message := &structs.Message{
		ID:         id,
		Name:       options.Name,
		Payload:    options.Payload,
		ResponseTo: options.ResponseTo,
		Sender:     e.Config.Role,
	}
	serializedBytes, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	return id, websocket.Message.Send(e.socket, serializedBytes)
}

// Helpers

func (e *ExoRelay) listenForMessages() {
	for {
		message, err := e.readMessage()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading message from websocket:", err)
			break
		} else {
			e.messageChannel <- message
		}
	}
	close(e.messageChannel)
}

func (e *ExoRelay) readMessage() (structs.Message, error) {
	var bytes []byte
	if err := websocket.Message.Receive(e.socket, &bytes); err != nil {
		return structs.Message{}, err
	}

	var unmarshaled structs.Message
	err := json.Unmarshal(bytes, &unmarshaled)
	if err != nil {
		return structs.Message{}, err
	}

	return unmarshaled, nil
}
