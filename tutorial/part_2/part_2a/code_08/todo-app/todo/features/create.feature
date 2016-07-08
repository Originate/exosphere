Feature: Creating todos

  Rules:
  - when successful, the service replies with "todo.created"
    and the created record
  - when there is an error, the service replies with "todo.not-created"
    and a message describing the error


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: creating a valid todo record
    When sending the message "todo.create" with the payload:
      """
      { "name": "Jean-Luc Picard" }
      """
    Then the service replies with "todo.created" and the payload:
      """
      {
        "id": /\d+/,
        "name": 'Jean-Luc Picard'
      }
      """
    And the service now contains the todos:
      | NAME            |
      | Jean-Luc Picard |
