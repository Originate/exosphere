package exoservice

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/Originate/exocom/go/exorelay"
	"github.com/Originate/exocom/go/structs"
)

// Request contains payload and helper functions for sending messages
type Request struct {
	Payload structs.MessagePayload
	Reply   func(exorelay.MessageOptions) error
	Send    func(exorelay.MessageOptions) error
}

// MessageHandler is the function signature for handling a message
type MessageHandler func(Request)

// MessageHandlerMapping is a map from message name to MessageHandler
type MessageHandlerMapping map[string]MessageHandler

// Bootstrap brings an ExoService online using environment variables and the
// given messageHandlers
func Bootstrap(messageHandlers MessageHandlerMapping) {
	port, err := strconv.Atoi(os.Getenv("EXOCOM_PORT"))
	if err != nil {
		panic(fmt.Sprintf("Invalid port: %s", os.Getenv("EXOCOM_PORT")))
	}
	config := exorelay.Config{
		Host: os.Getenv("EXOCOM_HOST"),
		Port: port,
		Role: os.Getenv("ROLE"),
	}
	service := ExoService{}
	err = service.Connect(config)
	if err != nil {
		panic(fmt.Sprintf("Error connecting: %v", err))
	}
	fmt.Println("online at port", port)
	err = service.ListenForMessages(messageHandlers)
	if err != nil {
		panic(fmt.Sprintf("Error listening for messages: %v", err))
	}
}

// ExoService is the high level Go API to talk to Exocom
type ExoService struct {
	exoRelay        exorelay.ExoRelay
	messageHandlers MessageHandlerMapping
}

// Connect brings an ExoService instance online
func (e *ExoService) Connect(config exorelay.Config) error {
	e.exoRelay = exorelay.ExoRelay{Config: config}
	return e.exoRelay.Connect()
}

// ListenForMessages blocks while the instance listens for and responds to messages
// returns when the message channel closes
func (e *ExoService) ListenForMessages(messageHandlers MessageHandlerMapping) error {
	if &e.exoRelay == nil {
		return errors.New("Please call Connect() first")
	}
	e.messageHandlers = messageHandlers
	messageChannel := e.exoRelay.GetMessageChannel()
	for {
		message, ok := <-messageChannel
		if !ok {
			return nil
		}
		e.receiveMessage(message)
	}
}

// Helpers

func (e *ExoService) receiveMessage(message structs.Message) {
	if e.messageHandlers == nil {
		return
	}
	if e.messageHandlers[message.Name] == nil {
		return
	}
	e.messageHandlers[message.Name](Request{
		Payload: message.Payload,
		Reply:   e.buildSendHelper(message.ID),
		Send:    e.buildSendHelper(""),
	})
}

func (e *ExoService) buildSendHelper(responseTo string) func(exorelay.MessageOptions) error {
	return func(options exorelay.MessageOptions) error {
		_, err := e.exoRelay.Send(exorelay.MessageOptions{
			Name:       options.Name,
			Payload:    options.Payload,
			ResponseTo: responseTo,
		})
		return err
	}
}
