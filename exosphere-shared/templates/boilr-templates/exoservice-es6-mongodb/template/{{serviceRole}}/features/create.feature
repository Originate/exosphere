Feature: Creating {{modelName}}s

  Rules:
  - when successful, the service replies with "{{modelName}}.created"
    and the created record
  - when there is an error, the service replies with "{{modelName}}.not-created"
    and a message describing the error


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: creating a valid {{modelName}} record
    When receiving the message "{{modelName}}.create" with the payload:
      """
      { "name": "one" }
      """
    Then the service replies with "{{modelName}}.created" and the payload:
      """
      {
        "id": /\d+/,
        "name": 'one'
      }
      """
    And the service now contains the {{modelName}}s:
      | NAME |
      | one  |
