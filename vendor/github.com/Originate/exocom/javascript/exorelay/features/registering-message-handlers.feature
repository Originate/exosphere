Feature: Registering message handlers

  As a deveoper writing Exosphere applications
  I want my application to respond to incoming messages
  So that I can write Exosphere services.

  Rules:
  - register message handlers by via "registerHandler"
  - the registered messages are invoked by sending a message to the ExoRelay socket
  - messages are executed asynchronously, and can send other messages back to indicate responses


  Background:
    Given ExoCom runs at port 4100
    And an ExoRelay instance called "exo-relay"


  Scenario: registering a message handler
    When I register a message handler:
      """
      exo-relay.register-handler 'message-1', ->
      """
    Then the instance has a handler for the message "message-1"


  Scenario: registering several handlers
    When I register the message handlers:
      """
      exo-relay.register-handlers message1: ->, message2: ->
      """
    Then the instance has handlers for the messages "message1" and "message2"



  # ERROR CHECKING

  Scenario: registering an already handled message
    Given my ExoRelay instance already has a handler for the message "hello"
    When I try to add another handler for that message
    Then ExoRelay emits an "error" event with the error "There is already a handler for message 'hello'"


  Scenario: forgetting to provide the message
    When I try to register the message handler:
      """
      exo-relay.register-handler ->
      """
    Then ExoRelay emits an "error" event with the error "Message ids must be strings"


  Scenario: providing an empty message
    When I try to register the message handler:
      """
      exo-relay.register-handler '', ->
      """
    Then ExoRelay emits an "error" event with the error "No message id provided"


  Scenario: providing a non-string message
    When I try to register the message handler:
      """
      exo-relay.register-handler [], ->
      """
    Then ExoRelay emits an "error" event with the error "Message ids must be strings"


  Scenario: forgetting to provide the handler
    When I try to register the message handler:
      """
      exo-relay.register-handler 'message'
      """
    Then ExoRelay emits an "error" event with the error "No message handler provided"


  Scenario: providing a non-functional handler
    When I try to register the message handler:
      """
      exo-relay.register-handler 'message', 'zonk'
      """
    Then ExoRelay emits an "error" event with the error "Message handler must be a function"
