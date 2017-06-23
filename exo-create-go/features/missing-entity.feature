Feature: Missing entity

  As a user unsure which entities I can create
  I want to be able to start the exo-create command without arguments and see the possible entities
  So that I can call this command correctly and build up muscle memory.

  - calling "exo-create" without an argument displays the list of entities that can be created


  Scenario: the user calls "exo create" without arguments
    When running "exo-create" in the terminal
    Then it prints "Error: missing entity for 'create' command" in the terminal
