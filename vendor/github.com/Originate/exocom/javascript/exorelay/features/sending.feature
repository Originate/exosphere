Feature: Sending outgoing messages

  As an Exoservice developer
  I want my service to be able to send messages to other Exosphere services
  So that it can interact with the rest of the application.

  Rules:
  - call "send" on your ExoRelay instance have it send out the given message
  - provide the message to send as the first parameter
  - provide payload for the message as the second parameter
  - the payload can be either a string, an array, or a Hash


  Background:
    Given ExoCom runs at port 4100
    And an ExoRelay instance called "exo-relay" running inside the "test-service" service


  Scenario: sending a message without payload
    When sending the message:
      """
      exo-relay.send 'hello-world'
      """
    Then ExoRelay makes the WebSocket request:
      """
      name: 'hello-world'
      sender: 'test-service'
      id: '<%= request_uuid %>'
      """


  Scenario: sending a message with a populated Hash as payload
    When sending the message:
      """
      exo-relay.send 'hello', name: 'world'
      """
    Then ExoRelay makes the WebSocket request:
      """
      name: 'hello'
      sender: 'test-service'
      payload:
        name: 'world'
      id: '<%= request_uuid %>'
      """


  Scenario: sending a message with an empty Hash as payload
    When sending the message:
      """
      exo-relay.send 'hello', {}
      """
    Then ExoRelay makes the WebSocket request:
      """
      name: 'hello'
      sender: 'test-service'
      payload: {}
      id: '<%= request_uuid %>'
      """


  Scenario: sending a message with a string as payload
    When sending the message:
      """
      exo-relay.send 'hello', 'world'
      """
    Then ExoRelay makes the WebSocket request:
      """
      name: 'hello'
      sender: 'test-service'
      payload: 'world'
      id: '<%= request_uuid %>'
      """

  Scenario: sending a message with an empty string as payload
    When sending the message:
      """
      exo-relay.send 'hello', ''
      """
    Then ExoRelay makes the WebSocket request:
      """
      name: 'hello'
      sender: 'test-service'
      payload: ''
      id: '<%= request_uuid %>'
      """


  Scenario: sending a message with an array as payload
    When sending the message:
      """
      exo-relay.send 'sum', [1, 2, 3]
      """
    Then ExoRelay makes the WebSocket request:
      """
      name: 'sum'
      sender: 'test-service'
      payload: [1, 2, 3]
      id: '<%= request_uuid %>'
      """


  Scenario: sending a message with an empty array as payload
    When sending the message:
      """
      exo-relay.send 'sum', []
      """
    Then ExoRelay makes the WebSocket request:
      """
      name: 'sum'
      sender: 'test-service'
      payload: []
      id: '<%= request_uuid %>'
      """


  Scenario: sending a message with whitespace in the name
    When sending the message:
      """
      exo-relay.send 'hello world'
      """
    Then ExoRelay makes the WebSocket request:
      """
      name: 'hello world'
      sender: 'test-service'
      id: '<%= request_uuid %>'
      """


  Scenario: trying to send an empty message
    When trying to send an empty message:
      """
      exo-relay.send ''
      """
    Then ExoRelay emits an "error" event with the error "ExoRelay#send cannot send empty messages"


  Scenario: trying to send a non-string message
    When trying to send a non-string message:
      """
      exo-relay.send []
      """
    Then ExoRelay emits an "error" event with the error "ExoRelay#send can only send string messages"


  Scenario: forgetting to provide the payload when providing a reply handler
    When trying to send a non-serializable payload:
      """
      exo-relay.send 'test', ->
      """
    Then ExoRelay emits an "error" event with the error "ExoRelay#send cannot send functions as payload"
