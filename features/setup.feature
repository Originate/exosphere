Feature: Setup of Exosphere applications

  As a developer starting to work on an existing Exosphere application
  I want to be able to set up the application to run
  So that I can start developing without having to research how to set up each individual service.

  - run "exo setup" to set up all the services in the current application


  Scenario: set up the "test" application
    Given a freshly checked out "test" application
    When running "exo setup" in this application's directory
    And waiting until the process ends
    Then it has created the folders:
      | SERVICE       | FOLDER       |
      | dashboard     | node_modules |
      | mongo-service | node_modules |
      | web-server    | node_modules |
