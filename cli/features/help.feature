Feature: executing abbreviated Exosphere commands

  As a developer using Exosphere
  I want to be able to use abbreviations for commands
  So that I can use Exosphere without having to type so much.

  Scenario: executing "exo help"
    When running "exo help" in the terminal
    Then it prints "Usage: exo <command>" in the terminal
    And it does not print "Error: " in the terminal

  Scenario Outline: executing "exo help <COMMAND>"
    When running "exo help <COMMAND>" in the terminal
    Then it prints "Usage: exo <COMMAND>" in the terminal
    And it does not print "Error: " in the terminal
  Examples:
    | COMMAND |
    | add     |
    | clean   |
    | clone   |
    | create  |
    | deploy  |
    | lint    |
    | run     |
    | setup   |
    | sync    |
    | test    |
