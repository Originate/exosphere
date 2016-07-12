Feature: Creating _____modelName_____s

  Rules:
  - when successful, the service replies with "_____modelName_____.created"
    and the created record
  - when there is an error, the service replies with "_____modelName_____.not-created"
    and a message describing the error


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: creating a valid _____modelName_____ record
    When receiving the message "_____modelName_____.create" with the payload:
      """
      { "name": "one" }
      """
    Then the service replies with "_____modelName_____.created" and the payload:
      """
      {
        "id": /\d+/,
        "name": 'one'
      }
      """
    And the service now contains the _____modelName_____s:
      | NAME |
      | one  |
