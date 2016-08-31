Feature: Updating a _____modelName_____

  Rules:
  - when receiving "_____modelName_____.update",
    updates the _____modelName_____ record with the given id
    and returns "_____modelName_____.updated" with the new record


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the _____modelName_____s:
      | NAME |
      | one  |
      | two  |


  Scenario: updating an existing _____modelName_____
    When receiving the message "_____modelName_____.update" with the payload:
      """
      {
        "id": "<%= idOf('one') %>",
        "name": "number one"
      }
      """
    Then the service replies with "_____modelName_____.updated" and the payload:
      """
      {
        "id": /.+/,
        "name": "number one"
      }
      """
    And the service now contains the _____modelName_____s:
      | NAME       |
      | number one |
      | two        |


  Scenario: trying to update a non-existing _____modelName_____
    When receiving the message "_____modelName_____.update" with the payload:
      """
      {
        "id": "zero",
        "name": "a total zero"
      }
      """
    Then the service replies with "_____modelName_____.not-found" and the payload:
      """
      {
        "id": "zero"
      }
      """
    And the service now contains the _____modelName_____s:
      | NAME |
      | one  |
      | two  |
