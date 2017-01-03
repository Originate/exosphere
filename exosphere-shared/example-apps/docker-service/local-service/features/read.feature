Feature: Get details for a test

  Rules:
  - when receiving "test.read",
    returns "test.details" with details for the given test


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the tests:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |


  Scenario: locating an existing test by id
    When sending the message "test.read" with the payload:
      """
      id: '<%= @id_of 'Jean-Luc Picard' %>'
      """
    Then the service replies with "test.details" and the payload:
      """
      id: /.+/
      name: 'Jean-Luc Picard'
      """


  Scenario: locating a non-existing test by id
    When sending the message "test.read" with the payload:
      """
      id: 'zonk'
      """
    Then the service replies with "test.not-found" and the payload:
      """
      id: 'zonk'
      """
