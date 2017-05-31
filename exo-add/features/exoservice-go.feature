Feature: scaffolding an ExoService written in Go

  As a Go developer adding features to an Exosphere application
  I want to have an easy way to scaffold an empty ExoService
  So that I don't waste time copy-and-pasting a bunch of code.

  - run "exo add service <name> exoservice-go"
    to add a new exo-service written in Go to the current application
  - run "exo add service" to add the service interactively


  Scenario: calling with all command line arguments
    Given I am in the root directory of an empty application called "test app"
    When running "exo-add service --service-role=users --service-type=users --author=test-author --template-name=exoservice-go --model=user --description=testing" in this application's directory
    Then my application contains the file "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0

      bus:
        type: exocom
        version: 0.21.7

      services:
        public:
        private:
          users:
            location: ./users
      """
    And my application contains the file "users/service.yml" with the content:
      """
      type: users
      description: testing
      author: test-author

      setup: glide install && go get github.com/DATA-DOG/godog
      startup:
        command: go run server.go
        online-text: online at port
      tests: go test

      messages:
        receives:
          - ping
        sends:
          - pong

      dependencies:
      """
    And my application contains the file "users/server.go" with the content:
      """
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
      """
    And my application contains the file "users/README.md" containing the text:
      """
      # USERS service
      > testing
      """
