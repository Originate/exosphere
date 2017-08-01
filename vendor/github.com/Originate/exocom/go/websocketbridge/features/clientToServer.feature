Feature:

  As an Exosphere developer
  I want to be able to send messages from a web client to my Exosphere application
  So that I can create a useful application

  Rules
  - If a message does not have a SessionId, a new one is generated for that user


  Background:
    Given a websocket bridge connected to Exocom


  Scenario: receiving a message without a SessionId
    When receiving this message from the client:
      """
      {
        "name": "hello",
        "id": "123"
      }
      """
    Then websocket bridge makes the websocket request:
      """
      {
        "name": "hello",
        "sessionId": "{{.sessionId}}",
        "sender": "websocket-test"
      }
      """


  Scenario: different clients get different sessionIds when they connnect
    When receiving a message "hello" from 2 different clients
    Then those clients have different sessionIds
