Feature: Listing all {{modelName}}s

  Rules:
  - returns all {{modelName}}s currently stored


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: no {{modelName}}s exist in the database
    When receiving the message "{{modelName}}.list"
    Then the service replies with "{{modelName}}.listing" and the payload:
      """
      []
      """


  Scenario: {{modelName}}s exist in the database
    Given the service contains the {{modelName}}s:
      | NAME |
      | one  |
      | two  |
    When receiving the message "{{modelName}}.list"
    Then the service replies with "{{modelName}}.listing" and the payload:
      """
      [
        {
          "name": "one",
          "id": /\d+/
        },
        {
          "name": "two",
          "id": /\d+/
        }
      ]
      """
