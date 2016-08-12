Feature: cloning an Exosphere application and its services

  As an application developer
  I want to download the complete source code for an application onto my machine
  So that I can work on it without manually cloning multiple repositories.

  - "exo clone" downloads the Exosphere application at the given URL onto the user's machine
  - the URL must point to a Git repository containing an Exosphere application,
    i.e. have an application.yml file in the root of the repository
  - all services referenced by the application.yml file also get cloned onto the machine


  Scenario: clone the "test" application with external dependencies.
    Given I am in an empty folder
    And source control contains the services "web" and "users"
    And source control contains a repo "app" with a file "application.yml" and the content:
      """
      name: Example application
      description: Demonstrates basic Exosphere application startup
      version: '1.0'

      services:
        web:
          local: ./web
          origin: ../origins/web
        users:
          local: ./users
          origin: ../origins/users
      """
    When running "exo clone origins/app" in the "tmp" directory
    Then it creates the files:
      | FILE            | FOLDER    |
      | application.yml | app       |
      | service.yml     | app/web   |
      | service.yml     | app/users |
