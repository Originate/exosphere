Feature: Handling incoming messages

  As a service developer
  I want to have an easy way to define code that handles incoming messages
  So that I can develop the services efficiently.


  Rules:
  - the handlers are defined in a file "server.ls"
  - this file exports a hash in which the key is the message name and the value the handler function
  - messages can be sent to the service as a POST request to "/run/<message-name>"
    with the data payload as the request body
  - data payload of messages goes in the request body, JSON encoded


  Background:
    Given an ExoCom instance
    And an instance of the "test service" service


  Scenario: receiving a message
    When receiving the "ping" message
    Then it acknowledges the received message
    And after a while it sends the "pong" message


  Scenario: receiving a message with payload
    When receiving the "greet" message with the payload:
      """
      name: 'world'
      """
    Then it acknowledges the received message
    And after a while it sends the "greeting" message with the textual payload:
      """
      Hello world
      """
