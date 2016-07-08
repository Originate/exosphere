Feature: Listing all _____serviceName_____s

  Rules:
  - returns all _____serviceName_____s currently stored


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: no _____serviceName_____s exist in the database
    When receiving the message "_____serviceName_____.list"
    Then the service replies with "_____serviceName_____s.listing" and the payload:
      """
      []
      """


  Scenario: _____serviceName_____s exist in the database
    Given the service contains the _____serviceName_____s:
      | NAME |
      | one  |
      | two  |
    When receiving the message "_____serviceName_____.list"
    Then the service replies with "_____serviceName_____.listing" and the payload:
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
