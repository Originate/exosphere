Feature: unknown entity

  As a user accidentally entering a non-existing entity
  I want to get a helpful error message and a list of entities that can be created
  So that I know how to use this command correctly.

  - calling "exo-create" with an unknown entity displays a help screen


  Scenario: the user tries to create an unknown entity
    When running "exo-create zonk" in the terminal
    Then it prints "Error: cannot create 'zonk'" in the terminal
