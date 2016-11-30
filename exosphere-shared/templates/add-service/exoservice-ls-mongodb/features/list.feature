Feature: Listing all _____modelName_____s

  Rules:
  - returns all _____modelName_____s currently stored


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: no _____modelName_____s exist in the database
    When sending the message "_____modelName_____.list"
    Then the service replies with "_____modelName_____s.listing" and the payload:
      """
      []
      """


  Scenario: _____modelName_____s exist in the database
    Given the service contains the _____modelName_____s:
      | NAME            |
      | Jean-Luc Picard |
      | Will Riker      |
    When sending the message "_____modelName_____.list"
    Then the service replies with "_____modelName_____.listing" and the payload:
      """
      [
        * name: 'Jean-Luc Picard'
          id: /\d+/
        * name: 'Will Riker'
          id: /\d+/
      ]
      """
