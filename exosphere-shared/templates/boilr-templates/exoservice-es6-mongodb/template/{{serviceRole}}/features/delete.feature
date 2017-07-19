Feature: Deleting a {{modelName}}

  Rules:
  - when receiving "{{modelName}}.delete",
    removes the {{modelName}} record with the given id
    and returns "{{modelName}}.deleted"


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the {{modelName}}s:
      | NAME |
      | one  |
      | two  |


  Scenario: deleting an existing {{modelName}}
    When receiving the message "{{modelName}}.delete" with the payload:
      """
      { "id": "<%= idOf('one') %>" }
      """
    Then the service replies with "{{modelName}}.deleted" and the payload:
      """
      {
        "id": /.+/,
        "name": 'one'
      }
      """
    And the service now contains the {{modelName}}s:
      | NAME |
      | two  |


  Scenario: trying to delete a non-existing {{modelName}}
    When receiving the message "{{modelName}}.delete" with the payload:
      """
      { "id": "zonk" }
      """
    Then the service replies with "{{modelName}}.not-found" and the payload:
      """
      { "id": "zonk" }
      """
    And the service now contains the {{modelName}}s:
      | NAME |
      | one  |
      | two  |
