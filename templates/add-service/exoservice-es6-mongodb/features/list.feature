Feature: Listing all _____serviceName_____s

  Rules:
  - returns all _____serviceName_____s currently stored


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: no _____serviceName_____s exist in the database
    When sending the message "_____serviceName_____.list"
    Then the service replies with "_____serviceName_____s.listing" and the payload:
      """
      []
      """


  Scenario: _____serviceName_____s exist in the database
    Given the service contains the _____serviceName_____s:
      | NAME            |
      | Jean-Luc Picard |
      | Will Riker      |
    When sending the message "_____serviceName_____.list"
    Then the service replies with "_____serviceName_____.listing" and the payload:
      """
      [
        {
          "name": "Jean-Luc Picard",
          "id": /\d+/
        },
        {
          "name": "Will Riker",
          "id": /\d+/
        }
      ]
      """
