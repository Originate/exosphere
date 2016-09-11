Feature: syncing Exosphere applications

  As an Exosphere developer
  I want to have an easy way to update my local clones of all code bases related to my project
  So that I work off the most current state of the product.

  Rules:
  - run "exo sync" in the directory of your application to sync it
  - this command runs "git pull" in all repositories denoted in the application.yml


  Scenario: syncing an application with external services
    Given a freshly checked out "external-service" application
    And The origin of "web-service" contains a new commit not yet present in the local clone
    When running "exo-sync" in this application's "app" directory
    Then it prints "Syncing application application with an external service" in the terminal
    And my application contains the newly committed file


  Scenario: syncing an application with only internal services
    Given a freshly checked out "simple" application
    When running "exo-sync" in this application's directory
    Then it prints "Sync successful" in the terminal


  Scenario: syncing an application with no services
    Given I am in the root directory of an empty application called "empty-app" with the file "application.yml":
      """
      name: empty-app
      description: A test application
      version: 0.0.1

      services:
      """
    When running "exo-sync" in this application's directory
    Then it prints "Sync successful" in the terminal
