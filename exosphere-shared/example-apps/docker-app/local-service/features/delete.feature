Feature: Deleting a test

  Rules:
  - when receiving "test.delete",
    removes the test record with the given id
    and returns "test.deleted"


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the tests:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |


  Scenario: deleting an existing test
    When sending the message "test.delete" with the payload:
      """
      id: '<%= @id_of 'Jean-Luc Picard' %>'
      """
    Then the service replies with "test.deleted" and the payload:
      """
      id: /.+/
      name: 'Jean-Luc Picard'
      """
    And the service now contains the tests:
      | NAME          |
      | William Riker |


  Scenario: trying to delete a non-existing test
    When sending the message "test.delete" with the payload:
      """
      id: 'zonk'
      """
    Then the service replies with "test.not-found" and the payload:
      """
      id: 'zonk'
      """
    And the service now contains the tests:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |
