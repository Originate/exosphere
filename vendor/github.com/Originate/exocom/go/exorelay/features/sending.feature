Feature: Sending outgoing messages

  As an Exoservice developer
  I want my service to be able to send messages to other Exosphere services
  So that it can interact with the rest of the application.

  Rules:
  - call "Send" on your ExoRelay instance have it send out the given message
  - provide the name of the message to send as the first parameter
  - provide payload for the message as the second parameter
  - the payload is a JSON structure
  - returns the message id which can be used to check if a received message is a reply


  Background:
    Given an ExoRelay with the role "test-service"
    And ExoRelay connects to Exocom


  Scenario: sending a message without payload
    When sending the message "hello-world"
    Then ExoRelay makes the WebSocket request:
      """
      {
        "name": "hello-world",
        "sender": "test-service",
        "id": "{{.outgoingMessageId}}"
      }
      """

  Scenario: sending a message with a populated Hash as payload
    When sending the message "hello" with the payload:
      """
      { "name": "world" }
      """
    Then ExoRelay makes the WebSocket request:
      """
      {
        "name": "hello",
        "payload": {
          "name": "world"
        },
        "sender": "test-service",
        "id": "{{.outgoingMessageId}}"
      }
      """

  Scenario: trying to send an empty message
    When trying to send an empty message
    Then ExoRelay errors with "ExoRelay#Send cannot send empty messages"
