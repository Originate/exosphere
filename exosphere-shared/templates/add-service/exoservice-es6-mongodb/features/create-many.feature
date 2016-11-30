Feature: Creating multiple _____modelName_____s

  As an ExoService application
  I want to be able to create multiple _____modelName_____ records in one transaction
  So that my application doesn't have to send and receive so many messages and remain performant.

  Rules:
  - send the message "_____modelName_____.create-many" to create several _____modelName_____ records at once
  - payload is an array of _____modelName_____ data
  - when successful, the service replies with "_____modelName_____.created-many"
    and the number of created records
  - when there is an error, the service replies with "_____modelName_____.not-created-many"
    and a message describing the error


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: creating valid _____modelName_____ records
    When receiving the message "_____modelName_____.create-many" with the payload:
      """
      [
        { "name": "one" },
        { "name": "two" }
      ]
      """
    Then the service replies with "_____modelName_____.created-many" and the payload:
      """
      { "count": 2 }
      """
    And the service now contains the _____modelName_____s:
      | NAME |
      | one  |
      | two  |
