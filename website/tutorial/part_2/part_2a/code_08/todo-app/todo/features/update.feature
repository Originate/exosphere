Feature: Updating a todo

  Rules:
  - when receiving "todo.update",
    updates the todo record with the given id
    and returns "todo.updated" with the new record


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the todos:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |


  Scenario: updating an existing todo
    When sending the message "todo.update" with the payload:
      """
      {
        "id": "<%= idOf('Jean-Luc Picard') %>",
        "name": "Cptn. Picard"
      }
      """
    Then the service replies with "todo.updated" and the payload:
      """
      {
        "id": /.+/,
        "name": "Cptn. Picard"
      }
      """
    And the service now contains the todos:
      | NAME          |
      | Cptn. Picard  |
      | William Riker |


  Scenario: trying to update a non-existing todo
    When sending the message "todo.update" with the payload:
      """
      {
        "id": "zonk",
        "name": "Cptn. Zonk"
      }
      """
    Then the service replies with "todo.not-found" and the payload:
      """
      {
        "id": "zonk"
      }
      """
    And the service now contains the todos:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |
