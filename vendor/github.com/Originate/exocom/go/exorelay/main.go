package exorelay

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Originate/exocom/go/structs"
	"github.com/Originate/exocom/go/utils"
	uuid "github.com/satori/go.uuid"

	"github.com/gorilla/websocket"
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
	socket, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%d", e.Config.Host, e.Config.Port), nil)
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

// Close takes this ExoRelay instance offline
func (e *ExoRelay) Close() error {
	err := e.socket.Close()
	e.socket = nil
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
	return id, e.socket.WriteMessage(websocket.TextMessage, serializedBytes)
}

// Helpers

func (e *ExoRelay) listenForMessages() {
	err := utils.ListenForMessages(e.socket, func(message structs.Message) {
		e.messageChannel <- message
	})
	close(e.messageChannel)
	if err != nil {
		fmt.Println("Exorelay error listening for messages", err)
	}
}
