Feature: Updating a _____modelName_____

  Rules:
  - when receiving "_____modelName_____.update",
    updates the _____modelName_____ record with the given id
    and returns "_____modelName_____.updated" with the new record


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the _____modelName_____s:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |


  Scenario: updating an existing _____modelName_____
    When sending the message "_____modelName_____.update" with the payload:
      """
      id: '<%= @id_of 'Jean-Luc Picard' %>'
      name: 'Cptn. Picard'
      """
    Then the service replies with "_____modelName_____.updated" and the payload:
      """
      id: /.+/
      name: 'Cptn. Picard'
      """
    And the service now contains the _____modelName_____s:
      | NAME          |
      | Cptn. Picard  |
      | William Riker |


  Scenario: trying to update a non-existing _____modelName_____
    When sending the message "_____modelName_____.update" with the payload:
      """
      id: 'zonk'
      name: 'Cptn. Zonk'
      """
    Then the service replies with "_____modelName_____.not-found" and the payload:
      """
      id: 'zonk'
      """
    And the service now contains the _____modelName_____s:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |
