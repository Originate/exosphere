Feature: syncing Exosphere applications

  As an Exosphere developer
  I want to have an easy way to update my local clones of all code bases related to my project
  So that I work off the most current state of the product.

  Rules:
  - run "exo sync" in the directory of your application to sync it
  - this command runs "git pull" in all repositories denoted in the application.yml


  @verbose
  Scenario: syncing an application with external services
    Given I am in the root directory of an application called "sync-test" using an external service "dashboard"
    And The origin of "dashboard" contains a new commit not yet present in the local clone
    When running "exo-sync" in this application's directory
    Then it prints "Syncing application sync-test" in the terminal
    And my application contains the newly committed file "dashboard/service.yml"


  Scenario: syncing an application with only internal services
    Given a freshly checked out "test" application
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
