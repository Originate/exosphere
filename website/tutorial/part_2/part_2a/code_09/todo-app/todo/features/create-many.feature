Feature: Creating multiple todos

  As an ExoService application
  I want to be able to create multiple todo records in one transaction
  So that my application doesn't have to send and receive so many messages and remain performant.

  Rules:
  - send the message "todo.create-many" to create several todo records at once
  - payload is an array of todo data
  - when successful, the service replies with "todo.created-many"
    and the number of created records
  - when there is an error, the service replies with "todo.not-created-many"
    and a message describing the error


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: creating valid todo records
    When receiving the message "todo.create-many" with the payload:
      """
      [
        { "name": "one" },
        { "name": "two" }
      ]
      """
    Then the service replies with "todo.created-many" and the payload:
      """
      { "count": 2 }
      """
    And the service now contains the todos:
      | NAME |
      | one  |
      | two  |
