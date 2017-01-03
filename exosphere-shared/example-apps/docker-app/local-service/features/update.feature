Feature: Updating a test

  Rules:
  - when receiving "test.update",
    updates the test record with the given id
    and returns "test.updated" with the new record


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the tests:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |


  Scenario: updating an existing test
    When sending the message "test.update" with the payload:
      """
      id: '<%= @id_of 'Jean-Luc Picard' %>'
      name: 'Cptn. Picard'
      """
    Then the service replies with "test.updated" and the payload:
      """
      id: /.+/
      name: 'Cptn. Picard'
      """
    And the service now contains the tests:
      | NAME          |
      | Cptn. Picard  |
      | William Riker |


  Scenario: trying to update a non-existing test
    When sending the message "test.update" with the payload:
      """
      id: 'zonk'
      name: 'Cptn. Zonk'
      """
    Then the service replies with "test.not-found" and the payload:
      """
      id: 'zonk'
      """
    And the service now contains the tests:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |
