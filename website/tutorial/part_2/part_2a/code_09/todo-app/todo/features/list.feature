Feature: Listing all todos

  Rules:
  - returns all todos currently stored


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: no todos exist in the database
    When receiving the message "todo.list"
    Then the service replies with "todo.listing" and the payload:
      """
      []
      """


  Scenario: todos exist in the database
    Given the service contains the todos:
      | NAME |
      | one  |
      | two  |
    When receiving the message "todo.list"
    Then the service replies with "todo.listing" and the payload:
      """
      [
        {
          "name": "one",
          "id": /\d+/
        },
        {
          "name": "two",
          "id": /\d+/
        }
      ]
      """
