Feature: Get details for a _____modelName_____

  Rules:
  - when receiving "_____modelName_____.read",
    returns "_____modelName_____.details" with details for the given _____modelName_____


  Background:
    Given an ExoCom server
    And an instance of this service
    And the service contains the _____modelName_____s:
      | NAME |
      | one  |
      | two  |


  Scenario: locating an existing _____modelName_____ by id
    When receiving the message "_____modelName_____.read" with the payload:
      """
      {
        "id": "<%= idOf('one') %>"
      }
      """
    Then the service replies with "_____modelName_____.details" and the payload:
      """
      {
        "id": /.+/,
        "name": "one"
      }
      """


  Scenario: locating a non-existing _____modelName_____ by id
    When receiving the message "_____modelName_____.read" with the payload:
      """
      {
        "id": "zonk"
      }
      """
    Then the service replies with "_____modelName_____.not-found" and the payload:
      """
      {
        "id": "zonk"
      }
      """
