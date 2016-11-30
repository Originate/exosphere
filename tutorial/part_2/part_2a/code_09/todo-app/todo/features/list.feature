Feature: Listing all todos

  Rules:
  - returns all todos currently stored


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: no todos exist in the database
    When sending the message "todo.list"
    Then the service replies with "todos.listing" and the payload:
      """
      []
      """


  Scenario: todos exist in the database
    Given the service contains the todos:
      | NAME            |
      | Jean-Luc Picard |
      | Will Riker      |
    When sending the message "todo.list"
    Then the service replies with "todo.listing" and the payload:
      """
      [
        {
          "name": "Jean-Luc Picard",
          "id": /\d+/
        },
        {
          "name": "Will Riker",
          "id": /\d+/
        }
      ]
      """
