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
        add         Adds a new service to the current application
        clean       Removes dangling Docker images and volumes
        configure   Configures secrets for an Exosphere application deployed to the cloud
        create      Creates a new Exosphere application
        deploy      Deploys Exosphere application to the cloud
        generate    Generates docker-compose and terraform files
        run         Runs an Exosphere application
        template    Manage service templates
        test        Runs tests for the application
        version     Displays the version

      Use "exo [command] --help" for more information about a command
      """


  Scenario: the user enters 'exo add help'
    When running "exo add help" in the terminal
    Then I see:
      """
      Adds a new service to the current application

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
      Creates a new Exosphere application

      Usage:
        exo create
      """

  Scenario: the user enters 'exo help deploy'
    When running "exo help deploy" in the terminal
    Then I see:
      """
      Deploys Exosphere application to the cloud

      Usage:
        exo deploy
      """

  Scenario: the user enters 'exo help generate'
    When running "exo help generate" in the terminal
    Then I see:
      """
      Generates docker-compose and terraform files

      Usage:
        exo generate
      """

  Scenario: the user enters 'exo template help'
    When running "exo template help" in the terminal
    Then I see:
      """
      Manage service templates

      Usage:
        exo template [command]

      Available Commands:
        test        Test a service template
      """

  Scenario: the user enters 'exo template test help'
    When running "exo template test help" in the terminal
    Then I see:
      """
      Test a service template by adding it to an exopshere application and running tests.
      This command must be run in the directory of an exosphere service template.

      Usage:
        exo template test
      """

  Scenario: the user enters 'exo run help'
    When running "exo run help" in the terminal
    Then I see:
      """
      Runs an Exosphere application

      Usage:
        exo run
      """
