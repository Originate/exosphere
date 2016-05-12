Feature: wrong command

  As an Exosphere user accidentally running the "exo" program without correct arguments
  I want to be given helpful information how to use this application
  So that I can do the things I intended without having to look this up somewhere else.


  Scenario: the user enters no command
    When trying to run "exo" in the terminal
    Then I see "Error: missing command"


  Scenario: the user enters an unknown command
    When trying to run "exo zonk" in the terminal
    Then I see "Error: unknown command 'zonk'"
