package main

import (
	"fmt"

	"github.com/Originate/exocom/go/exorelay"
	"github.com/Originate/exocom/go/exoservice"
)

func handlePing(request exoservice.Request) {
	err := request.Reply(exorelay.MessageOptions{Name: "pong"})
	if err != nil {
		panic(fmt.Sprintf("Failed to send reply: %v", err))
	}
}

func main() {
	messageHandlers := exoservice.MessageHandlerMapping{
		"ping": handlePing,
	}
	exoservice.Bootstrap(messageHandlers)
}
