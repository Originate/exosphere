Feature: Installing Exosphere applications

  As a developer starting to work on an existing Exosphere application
  I want to be able to set up the application to run
  So that I can start developing without having to research how to install each individual service.

  - run "exo install" to set up all the services in the current application


  Scenario: installing the "test" application
    Given a freshly checked out "test" application
    When installing it
    Then it creates the folders:
      | SERVICE       | FOLDER       |
      | dashboard     | node_modules |
      | mongo-service | node_modules |
      | web-server    | node_modules |
