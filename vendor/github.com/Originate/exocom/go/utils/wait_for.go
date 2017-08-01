package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/Originate/exocom/go/structs"
)

// WaitFor waits a maximum of 10 seconds the given condition to become true
// returning an error with the given message if and only if it does not
func WaitFor(condition func() bool, errorMessage string) error {
	return WaitForf(condition, func() error {
		return errors.New(errorMessage)
	})
}

// WaitForf is similar to WaitFor but accepts a func to generate the error message
func WaitForf(condition func() bool, errorFn func() error) error {
	success := make(chan bool, 1)

	go func() {
		for {
			time.Sleep(time.Millisecond)
			if condition() {
				success <- true
				break
			}
		}
	}()

	select {
	case <-success:
		return nil
	case <-time.After(time.Second * 10):
		return errorFn()
	}
}

// ObjectWithMessages is the interface for a object passed into WaitForMessageWithName
type ObjectWithMessages interface {
	GetReceivedMessages() []structs.Message
}

// WaitForMessageWithName waits to receive a message with the given name
func WaitForMessageWithName(obj ObjectWithMessages, name string) (structs.Message, error) {
	var savedMessage structs.Message
	err := WaitFor(func() bool {
		for _, message := range obj.GetReceivedMessages() {
			if message.Name == name {
				savedMessage = message
				return true
			}
		}
		return false
	}, fmt.Sprintf("Expected to recieve a message with name '%s' but only has %v", name, obj.GetReceivedMessages()))
	return savedMessage, err
}
