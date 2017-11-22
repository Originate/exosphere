Feature: Get details for a todo

  Rules:
  - when receiving "todo.read",
    returns "todo.details" with details for the given todo


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the todos:
      | NAME |
      | one  |
      | two  |


  Scenario: locating an existing todo by id
    When receiving the message "todo.read" with the payload:
      """
      {
        "id": "<%= idOf('one') %>"
      }
      """
    Then the service replies with "todo.details" and the payload:
      """
      {
        "id": /.+/,
        "name": "one"
      }
      """


  Scenario: locating a non-existing todo by id
    When receiving the message "todo.read" with the payload:
      """
      {
        "id": "zonk"
      }
      """
    Then the service replies with "todo.not-found" and the payload:
      """
      {
        "id": "zonk"
      }
      """
