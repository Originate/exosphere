Feature: Get details for a _____modelName_____

  Rules:
  - when receiving "_____modelName_____.read",
    returns "_____modelName_____.details" with details for the given _____modelName_____


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the _____modelName_____s:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |


  Scenario: locating an existing _____modelName_____ by id
    When sending the message "_____modelName_____.read" with the payload:
      """
      id: '<%= @id_of 'Jean-Luc Picard' %>'
      """
    Then the service replies with "_____modelName_____.details" and the payload:
      """
      id: /.+/
      name: 'Jean-Luc Picard'
      """


  Scenario: locating a non-existing _____modelName_____ by id
    When sending the message "_____modelName_____.read" with the payload:
      """
      id: 'zonk'
      """
    Then the service replies with "_____modelName_____.not-found" and the payload:
      """
      id: 'zonk'
      """
