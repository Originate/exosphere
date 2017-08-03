Feature: 

  As an Exosphere developer
  I want to be able to send messages from my Exosphere application to a client
  So that I can create a useful application

  Rules
  - messages received from exorelay are forwarded to a client only if the session id matches


  Background:
    Given a websocket bridge connected to Exocom


  Scenario: receiving a message with a SessionId
    When receiving this message from the client:
      """
      {
        "name": "ping",
        "id": "123"
      }
      """
    And receiving this message from exocom:
      """
      {
        "name": "pong",
        "sessionId": "new-session-id",
        "id": "456"
      }
      """
    Then websocket bridge sends this to the client:
      """
      {
        "name": "pong",
        "id": "456"
      }
      """


  Scenario: receiving a message with no SessionId
    When receiving this message from the client:
      """
      {
        "name": "ping",
        "id": "123"
      }
      """
    And receiving this message from exocom:
      """
      {
        "name": "pong",
        "id": "456"
      }
      """
    Then the websocket bridge does not send a message to the client
