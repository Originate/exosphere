package exoserviceTestFixtures

import (
	"fmt"

	"github.com/Originate/exocom/go/exorelay"
	"github.com/Originate/exocom/go/exoservice"
	"github.com/Originate/exocom/go/structs"
)

// PingTextFixture is a test fixture which responds to "ping" messages with "pong" messages
type PingTextFixture struct {
	ReceivedMessages []structs.Message
}

// GetMessageHandler returns a message hangler
func (r *PingTextFixture) GetMessageHandler() exoservice.MessageHandlerMapping {
	return exoservice.MessageHandlerMapping{
		"ping": func(request exoservice.Request) {
			err := request.Reply(exorelay.MessageOptions{Name: "pong"})
			if err != nil {
				panic(fmt.Sprintf("Failed to send reply: %v", err))
			}
		},
		"ping it": func(request exoservice.Request) {
			err := request.Reply(exorelay.MessageOptions{Name: "pong it"})
			if err != nil {
				panic(fmt.Sprintf("Failed to send reply: %v", err))
			}
		},
	}
}
