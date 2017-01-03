Feature: Creating multiple tests

  As an ExoService application
  I want to be able to create multiple test records in one transaction
  So that my application doesn't have to send and receive so many messages and remain performant.

  Rules:
  - send the message "test.create-many" to create several test records at once
  - payload is an array of test data
  - when successful, the service replies with "test.created-many"
    and the number of newly created records
  - when there is an error, the service replies with "test.not-created-many"
    and a message describing the error


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: creating valid test records
    When sending the message "test.create-many" with the payload:
      """
      [
        * name: 'Jean-Luc Picard'
        * name: 'William Riker'
      ]
      """
    Then the service replies with "test.created-many" and the payload:
      """
      count: 2
      """
    And the service now contains the tests:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |
