Feature: Deleting a todo

  Rules:
  - when receiving "todo.delete",
    removes the todo record with the given id
    and returns "todo.deleted"


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the todos:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |


  Scenario: deleting an existing todo
    When sending the message "todo.delete" with the payload:
      """
      { "id": "<%= idOf('Jean-Luc Picard') %>" }
      """
    Then the service replies with "todo.deleted" and the payload:
      """
      {
        "id": /.+/,
        "name": 'Jean-Luc Picard'
      }
      """
    And the service now contains the todos:
      | NAME          |
      | William Riker |


  Scenario: trying to delete a non-existing todo
    When sending the message "todo.delete" with the payload:
      """
      { "id": "zonk" }
      """
    Then the service replies with "todo.not-found" and the payload:
      """
      { "id": "zonk" }
      """
    And the service now contains the todos:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |
