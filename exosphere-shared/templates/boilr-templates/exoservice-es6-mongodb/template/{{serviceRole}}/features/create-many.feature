Feature: Creating multiple {{modelName}}s

  As an ExoService application
  I want to be able to create multiple {{modelName}} records in one transaction
  So that my application doesn't have to send and receive so many messages and remain performant.

  Rules:
  - send the message "{{modelName}}.create-many" to create several {{modelName}} records at once
  - payload is an array of {{modelName}} data
  - when successful, the service replies with "{{modelName}}.created-many"
    and the number of created records
  - when there is an error, the service replies with "{{modelName}}.not-created-many"
    and a message describing the error


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: creating valid {{modelName}} records
    When receiving the message "{{modelName}}.create-many" with the payload:
      """
      [
        { "name": "one" },
        { "name": "two" }
      ]
      """
    Then the service replies with "{{modelName}}.created-many" and the payload:
      """
      { "count": 2 }
      """
    And the service now contains the {{modelName}}s:
      | NAME |
      | one  |
      | two  |
