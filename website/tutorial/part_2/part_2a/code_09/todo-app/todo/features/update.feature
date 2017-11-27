Feature: Updating a todo

  Rules:
  - when receiving "todo.update",
    updates the todo record with the given id
    and returns "todo.updated" with the new record


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the todos:
      | NAME |
      | one  |
      | two  |


  Scenario: updating an existing todo
    When receiving the message "todo.update" with the payload:
      """
      {
        "id": "<%= idOf('one') %>",
        "name": "number one"
      }
      """
    Then the service replies with "todo.updated" and the payload:
      """
      {
        "id": /.+/,
        "name": "number one"
      }
      """
    And the service now contains the todos:
      | NAME       |
      | number one |
      | two        |


  Scenario: trying to update a non-existing todo
    When receiving the message "todo.update" with the payload:
      """
      {
        "id": "zero",
        "name": "a total zero"
      }
      """
    Then the service replies with "todo.not-found" and the payload:
      """
      {
        "id": "zero"
      }
      """
    And the service now contains the todos:
      | NAME |
      | one  |
      | two  |
