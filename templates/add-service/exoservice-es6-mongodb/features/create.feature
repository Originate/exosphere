Feature: Creating _____serviceName_____s

  Rules:
  - when successful, the service replies with "_____serviceName_____.created"
    and the created record
  - when there is an error, the service replies with "_____serviceName_____.not-created"
    and a message describing the error


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: creating a valid _____serviceName_____ record
    When sending the message "_____serviceName_____.create" with the payload:
      """
      { "name": "Jean-Luc Picard" }
      """
    Then the service replies with "_____serviceName_____.created" and the payload:
      """
      {
        "id": /\d+/,
        "name": 'Jean-Luc Picard'
      }
      """
    And the service now contains the _____serviceName_____s:
      | NAME            |
      | Jean-Luc Picard |
