Feature: Updating a _____serviceName_____

  Rules:
  - when receiving "_____serviceName_____.update",
    updates the _____serviceName_____ record with the given id
    and returns "_____serviceName_____.updated" with the new record


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the _____serviceName_____s:
      | NAME |
      | one  |
      | two  |


  Scenario: updating an existing _____serviceName_____
    When receiving the message "_____serviceName_____.update" with the payload:
      """
      {
        "id": "<%= idOf('one') %>",
        "name": "number one"
      }
      """
    Then the service replies with "_____serviceName_____.updated" and the payload:
      """
      {
        "id": /.+/,
        "name": "number one"
      }
      """
    And the service now contains the _____serviceName_____s:
      | NAME       |
      | number one |
      | two        |


  Scenario: trying to update a non-existing _____serviceName_____
    When receiving the message "_____serviceName_____.update" with the payload:
      """
      {
        "id": "zero",
        "name": "a total zero"
      }
      """
    Then the service replies with "_____serviceName_____.not-found" and the payload:
      """
      {
        "id": "zero"
      }
      """
    And the service now contains the _____serviceName_____s:
      | NAME |
      | one  |
      | two  |
