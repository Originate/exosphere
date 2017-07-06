Feature: help command

  As a developer not familiar with how to call "exo clean" correctly
  I want to be given helpful information how to use this application
  So that I can do the things I intended without having to look this up somewhere else.


  Scenario: the user enters 'exo help'
    When running "exo help" in the terminal
    Then I see:
      """
      Usage:
        exo [command]

      Available Commands:
        add             Adds a new service to the current application
        clean           Removes dangling Docker images and volumes
        create          Create a new Exosphere application
        fetch-templates Fetch service templates
        version         Exosphere go version number

      Use "exo [command] --help" for more information about a command
      """


  Scenario: the user enters 'exo add help'
    When running "exo add help" in the terminal
    Then I see:
      """
      Adds a new service to the current application
      This command must be called in the root directory of the application

      Usage:
        exo add
      """


  Scenario: the user enters 'exo help clean'
    When running "exo clean help" in the terminal
    Then I see:
      """
      Removes dangling Docker images and volumes

      Usage:
        exo clean
      """


  Scenario: the user enters 'exo create help'
    When running "exo create help" in the terminal
    Then I see:
      """
      Create a new Exosphere application

      Usage:
        exo create
      """

  Scenario: the user enters 'exo fetch-templates help'
    When running "exo fetch-templates help" in the terminal
    Then I see:
      """
      Fetch service templates
      This command must be called in the root directory of the application

      Usage:
        exo fetch-templates
      """
