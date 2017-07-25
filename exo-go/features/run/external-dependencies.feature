Feature: running Exosphere applications with external dependencies

  As an Exosphere developer
  I want to have an easy way to run applications with external dependencies
  So that I can test my app locally.

  Rules:
  - run "exo run" in the directory of your application to run it
  - this command boots up all the services and dependencies of the application


  Scenario: booting an Exosphere application with external dependencies
    Given I am in the root directory of the "external-dependency" example application
    When starting "exo run" in my application directory
    Then it prints "setup complete" in the terminal
    And my machine has acquired the Docker images:
      | tmp_mongo |
      | mongo     |
    And the docker images have the following folders:
      | IMAGE       | FOLDER       |
      | tmp_mongo   | node_modules |
    And it prints "all dependencies online" in the terminal
    And it prints "all services online" in the terminal
    And my machine is running the services:
      | NAME  |
      | mongo |
