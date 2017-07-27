package exorelayTestFixtures

import (
	"fmt"

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

// WaitForReceivedMessagesCount waits the received messages count to equal the given count
func (r *ReceivingMessagesTestFixture) WaitForReceivedMessagesCount(count int) error {
	return utils.WaitFor(func() bool {
		return len(r.ReceivedMessages) >= count
	}, fmt.Sprintf("Expected to recieve %d messages but only has %d:\n%v", count, len(r.ReceivedMessages), r.ReceivedMessages))
}
