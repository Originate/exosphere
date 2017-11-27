Feature: Deleting a todo

  Rules:
  - when receiving "todo.delete",
    removes the todo record with the given id
    and returns "todo.deleted"


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the todos:
      | NAME |
      | one  |
      | two  |


  Scenario: deleting an existing todo
    When receiving the message "todo.delete" with the payload:
      """
      { "id": "<%= idOf('one') %>" }
      """
    Then the service replies with "todo.deleted" and the payload:
      """
      {
        "id": /.+/,
        "name": 'one'
      }
      """
    And the service now contains the todos:
      | NAME |
      | two  |


  Scenario: trying to delete a non-existing todo
    When receiving the message "todo.delete" with the payload:
      """
      { "id": "zonk" }
      """
    Then the service replies with "todo.not-found" and the payload:
      """
      { "id": "zonk" }
      """
    And the service now contains the todos:
      | NAME |
      | one  |
      | two  |
