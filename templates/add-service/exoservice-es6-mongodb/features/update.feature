Feature: Updating a _____serviceName_____

  Rules:
  - when receiving "_____serviceName_____.update",
    updates the _____serviceName_____ record with the given id
    and returns "_____serviceName_____.updated" with the new record


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the _____serviceName_____s:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |


  Scenario: updating an existing _____serviceName_____
    When sending the message "_____serviceName_____.update" with the payload:
      """
      {
        "id": "<%= idOf('Jean-Luc Picard') %>",
        "name": "Cptn. Picard"
      }
      """
    Then the service replies with "_____serviceName_____.updated" and the payload:
      """
      {
        "id": /.+/,
        "name": "Cptn. Picard"
      }
      """
    And the service now contains the _____serviceName_____s:
      | NAME          |
      | Cptn. Picard  |
      | William Riker |


  Scenario: trying to update a non-existing _____serviceName_____
    When sending the message "_____serviceName_____.update" with the payload:
      """
      {
        "id": "zonk",
        "name": "Cptn. Zonk"
      }
      """
    Then the service replies with "_____serviceName_____.not-found" and the payload:
      """
      {
        "id": "zonk"
      }
      """
    And the service now contains the _____serviceName_____s:
      | NAME            |
      | Jean-Luc Picard |
      | William Riker   |
