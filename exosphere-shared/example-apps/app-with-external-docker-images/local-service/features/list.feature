Feature: Listing all tests

  Rules:
  - returns all tests currently stored


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: no tests exist in the database
    When sending the message "test.list"
    Then the service replies with "test.listing" and the payload:
      """
      []
      """


  Scenario: tests exist in the database
    Given the service contains the tests:
      | NAME            |
      | Jean-Luc Picard |
      | Will Riker      |
    When sending the message "test.list"
    Then the service replies with "test.listing" and the payload:
      """
      [
        * name: 'Jean-Luc Picard'
          id: /\d+/
        * name: 'Will Riker'
          id: /\d+/
      ]
      """
