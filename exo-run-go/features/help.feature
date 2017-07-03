Feature: help command

  As a developer not familiar with how to call "exo run" correctly
  I want to be given helpful information how to use this application
  So that I can do the things I intended without having to look this up somewhere else.


  Scenario: the user enters 'exo-run help'
    When running "exo-run help" in the terminal
    Then I see:
    """
    Usage: exo run

    Runs an Exosphere application.
    This command must be run in the root directory of the application.
    """
