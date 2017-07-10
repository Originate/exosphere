Feature: Get details for a {{modelName}}

  Rules:
  - when receiving "{{modelName}}.read",
    returns "{{modelName}}.details" with details for the given {{modelName}}


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the {{modelName}}s:
      | NAME |
      | one  |
      | two  |


  Scenario: locating an existing {{modelName}} by id
    When receiving the message "{{modelName}}.read" with the payload:
      """
      {
        "id": "<%= idOf('one') %>"
      }
      """
    Then the service replies with "{{modelName}}.details" and the payload:
      """
      {
        "id": /.+/,
        "name": "one"
      }
      """


  Scenario: locating a non-existing {{modelName}} by id
    When receiving the message "{{modelName}}.read" with the payload:
      """
      {
        "id": "zonk"
      }
      """
    Then the service replies with "{{modelName}}.not-found" and the payload:
      """
      {
        "id": "zonk"
      }
      """
