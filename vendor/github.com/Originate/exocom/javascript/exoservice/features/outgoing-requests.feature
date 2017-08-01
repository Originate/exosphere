Feature: Sending outgoing requests

  As an Exoservice developer
  I want to be able to send other messages from my services
  So that my services can use other services.

  Rules:
  - use the "send" method given to your message handlers as a parameter


  Background:
    Given an ExoCom instance
    And an instance of the "test service" service


  Scenario: a service sends out a message
    When receiving the "sender" message
    Then after a while it sends the "greetings" message with the textual payload:
      """
      from the sender service
      """


  Scenario: a service sends out a message with whitespace
    When receiving the "some salutation" message
    Then after a while it sends the "this salutation" message with the textual payload:
      """
      salutations
      """
