package exoserviceTestFixtures

import (
	"fmt"

	"github.com/Originate/exocom/go/exorelay"
	"github.com/Originate/exocom/go/exoservice"
	"github.com/Originate/exocom/go/structs"
)

// SendTextFixture is a test fixture which sends "ping_received" messages when it receives "ping" messages
type SendTextFixture struct {
	ReceivedMessages []structs.Message
}

// GetMessageHandler returns a message hangler
func (r *SendTextFixture) GetMessageHandler() exoservice.MessageHandlerMapping {
	return exoservice.MessageHandlerMapping{
		"ping": func(request exoservice.Request) {
			err := request.Send(exorelay.MessageOptions{Name: "pong"})
			if err != nil {
				panic(fmt.Sprintf("Failed to send message: %v", err))
			}
		},
	}
}
