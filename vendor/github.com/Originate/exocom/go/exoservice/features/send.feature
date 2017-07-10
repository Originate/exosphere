Feature: Outgoing messages

  As an Exoservice developer
  I want to be able to send other messages from my services
  So that my services can use other services.

  Rules:
  - use the "helpers.Send" method given as a parameter to your message handler
    to send replies for the message you are currently handling


  Background:
    Given I connect the "send" test fixture


  Scenario: replying to a message
    When receiving a "ping" message
    Then it sends a "pong" message


  Scenario: replying to a message with whitespace
    When receiving a "ping it" message
    Then it sends a "pong it" message
