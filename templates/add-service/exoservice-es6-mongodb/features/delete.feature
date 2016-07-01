Feature: Deleting a _____serviceName_____

  Rules:
  - when receiving "_____serviceName_____.delete",
    removes the _____serviceName_____ record with the given id
    and returns "_____serviceName_____.deleted"


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the _____serviceName_____s:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |


  Scenario: deleting an existing _____serviceName_____
    When sending the message "_____serviceName_____.delete" with the payload:
      """
      { "id": "<%= idOf('Jean-Luc Picard') %>" }
      """
    Then the service replies with "_____serviceName_____.deleted" and the payload:
      """
      {
        "id": /.+/,
        "name": 'Jean-Luc Picard'
      }
      """
    And the service now contains the _____serviceName_____s:
      | NAME          |
      | William Riker |


  Scenario: trying to delete a non-existing _____serviceName_____
    When sending the message "_____serviceName_____.delete" with the payload:
      """
      { "id": "zonk" }
      """
    Then the service replies with "_____serviceName_____.not-found" and the payload:
      """
      { "id": "zonk" }
      """
    And the service now contains the _____serviceName_____s:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |
