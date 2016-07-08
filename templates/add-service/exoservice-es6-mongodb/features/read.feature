Feature: Get details for a _____serviceName_____

  Rules:
  - when receiving "_____serviceName_____.read",
    returns "_____serviceName_____.details" with details for the given _____serviceName_____


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the _____serviceName_____s:
      | NAME |
      | one  |
      | two  |


  Scenario: locating an existing _____serviceName_____ by id
    When receiving the message "_____serviceName_____.read" with the payload:
      """
      {
        "id": "<%= idOf('one') %>"
      }
      """
    Then the service replies with "_____serviceName_____.details" and the payload:
      """
      {
        "id": /.+/,
        "name": "one"
      }
      """


  Scenario: locating a non-existing _____serviceName_____ by id
    When receiving the message "_____serviceName_____.read" with the payload:
      """
      {
        "id": "zonk"
      }
      """
    Then the service replies with "_____serviceName_____.not-found" and the payload:
      """
      {
        "id": "zonk"
      }
      """
