Feature: Creating multiple _____serviceName_____s

  As an ExoService application
  I want to be able to create multiple _____serviceName_____ records in one transaction
  So that my application doesn't have to send and receive so many messages and remain performant.

  Rules:
  - send the message "_____serviceName_____.create-many" to create several _____serviceName_____ records at once
  - payload is an array of _____serviceName_____ data
  - when successful, the service replies with "_____serviceName_____.created-many"
    and the number of newly created records
  - when there is an error, the service replies with "_____serviceName_____.not-created-many"
    and a message describing the error


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: creating valid _____serviceName_____ records
    When sending the message "_____serviceName_____.create-many" with the payload:
      """
      [
        * name: 'Jean-Luc Picard'
        * name: 'William Riker'
      ]
      """
    Then the service replies with "_____serviceName_____.created-many" and the payload:
      """
      count: 2
      """
    And the service now contains the _____serviceName_____s:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |
