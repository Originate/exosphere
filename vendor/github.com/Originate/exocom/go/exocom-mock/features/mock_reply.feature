Feature: Mock replies

  Background:
    Given an Exocom server
    And I have mocked it to reply to "request message" with "reply message" and the payload:
      """
      {
        "some": "data"
      }
      """

  Scenario: receiving a message for which a mock reply is registered
    When receiving a "request message" message with id "123"
    Then it sends a "reply message" message as a reply "123" with the payload
      """
      {
        "some": "data"
      }
      """

  Scenario: receiving a message for which no mock reply is registered
    When receiving a "other message" message with id "123"
    Then it does not send a message
