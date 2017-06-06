Feature: help command

  As a developer not familiar with how to call "exo clean" correctly
  I want to be given helpful information how to use this application
  So that I can do the things I intended without having to look this up somewhere else.


  Scenario: the user enters 'exo-clean help'
    When running "exo-clean help" in the terminal
    Then I see:
    """
    Usage: exo clean

    Removes dangling Docker images and volumes.
    """
