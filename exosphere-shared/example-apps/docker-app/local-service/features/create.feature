Feature: Creating tests

  Rules:
  - when successful, the service replies with "test.created"
    and the newly created record
  - when there is an error, the service replies with "test.not-created"
    and a message describing the error


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: creating a valid test account
    When sending the message "test.create" with the payload:
      """
      name: 'Jean-Luc Picard'
      """
    Then the service replies with "test.created" and the payload:
      """
      id: /\d+/
      name: 'Jean-Luc Picard'
      """
    And the service now contains the tests:
      | NAME            |
      | Jean-Luc Picard |
