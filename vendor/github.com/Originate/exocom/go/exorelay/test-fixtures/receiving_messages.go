package exorelayTestFixtures

import (
	"github.com/Originate/exocom/go/exorelay"
	"github.com/Originate/exocom/go/structs"
	"github.com/Originate/exocom/go/utils"
)

// ReceivingMessagesTestFixture is a test fixture which saves the messages it receives
type ReceivingMessagesTestFixture struct {
	ReceivedMessages []structs.Message
}

// Setup setups up the test fixture for the given exorelay instance
func (r *ReceivingMessagesTestFixture) Setup(exoRelay *exorelay.ExoRelay) {
	messageChannel := exoRelay.GetMessageChannel()
	go func() {
		for {
			message, ok := <-messageChannel
			if !ok {
				break // channel closed
			}
			r.ReceivedMessages = append(r.ReceivedMessages, message)
		}
	}()
}

// GetReceivedMessages return the received messages
func (r *ReceivingMessagesTestFixture) GetReceivedMessages() []structs.Message {
	return r.ReceivedMessages
}

// WaitForMessageWithName waits to receive a message with the given name
func (r *ReceivingMessagesTestFixture) WaitForMessageWithName(name string) (structs.Message, error) {
	return utils.WaitForMessageWithName(r, name)
}
