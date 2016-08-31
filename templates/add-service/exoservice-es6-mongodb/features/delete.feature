Feature: Deleting a _____modelName_____

  Rules:
  - when receiving "_____modelName_____.delete",
    removes the _____modelName_____ record with the given id
    and returns "_____modelName_____.deleted"


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the _____modelName_____s:
      | NAME |
      | one  |
      | two  |


  Scenario: deleting an existing _____modelName_____
    When receiving the message "_____modelName_____.delete" with the payload:
      """
      { "id": "<%= idOf('one') %>" }
      """
    Then the service replies with "_____modelName_____.deleted" and the payload:
      """
      {
        "id": /.+/,
        "name": 'one'
      }
      """
    And the service now contains the _____modelName_____s:
      | NAME |
      | two  |


  Scenario: trying to delete a non-existing _____modelName_____
    When receiving the message "_____modelName_____.delete" with the payload:
      """
      { "id": "zonk" }
      """
    Then the service replies with "_____modelName_____.not-found" and the payload:
      """
      { "id": "zonk" }
      """
    And the service now contains the _____modelName_____s:
      | NAME |
      | one  |
      | two  |
