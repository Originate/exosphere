Feature: Sending requests to services

  As a developer TDD-ing an ExoService
  I want to be able to send it messages in my tests
  So that I can trigger desired activities in my service and observe its behavior.

  Rules:
  - you must register services with the mock exocomm instance before you can send messages to them
  - call "exocomm.sendMessage service: <service-name>, name: <message name>, payload: <payload>"
    to send the given message to the given service


  Background:
    Given an ExoComMock instance
    And a known "users" service


  Scenario: sending a message to a registered service
    When sending a "users.create" message to the "users" service with the payload:
      """
      name: 'Jean-Luc Picard'
      """
    Then ExoComMock makes the request:
      | NAME    | users.create            |
      | PAYLOAD | name: 'Jean-Luc Picard' |


  Scenario: sending a message to an unknown service
    When trying to send a "users.create" message to the "zonk" service
    Then I get the error "unknown service: 'zonk'"
