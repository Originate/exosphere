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
    And an ExoRelay instance running inside the "test-service" service at port 4000


  Scenario: sending a message without payload
    When sending the message:
      """
      { "name": "hello-world" }
      """
    Then ExoRelay makes the ZMQ request:
      """
      {
      "name": "hello-world",
      "sender": "test-service",
      "id": "<request_uuid>"
      }
      """


  Scenario: sending a message with a JSON object (string) as payload
    When sending the message:
      """
      {
      "name": "hello",
        "payload": { "name" : "world" }
      }
      """
    Then ExoRelay makes the ZMQ request:
      """
      {
      "name": "hello",
      "sender": "test-service",
      "payload": { "name" : "world" },
      "id": "<request_uuid>"
      }
      """



  Scenario: sending a message with a JSON object (array) as payload
    When sending the message:
      """
      {
      "name": "sum",
      "payload": { "array" : [1,2,3] }
      }
      """
    Then ExoRelay makes the ZMQ request:
      """
      {
      "name": "sum",
      "sender": "test-service",
      "payload": { "array" : [1,2,3] },
      "id": "<request_uuid>"
      }
      """


  Scenario: trying to send an empty message
    When trying to send an empty message:
      """
      {}
      """
    Then ExoRelay throws an ExorelayError exception with the message "ExoRelay#send cannot send empty messages"

