Feature: running Exosphere applications

  As an Exosphere developer
  I want to have an easy way to run an application on my developer machine
  So that I can test my app locally.

  Rules:
  - run "exo run" in the directory of your application to run it
  - this command boots up all the services and dependencies of the application


  Scenario: booting a functioning Exosphere application
    Given I am in the root directory of the "running" example application
    And my application has been set up correctly
    When starting "exo run" in my application directory
    Then my machine is running the services:
      | NAME  |
      | web   |
    And it prints "all services online" in the terminal


  Scenario: booting an Exosphere application with external docker images
    Given I am in the root directory of the "app-with-external-docker-images" example application
    And my application has been set up correctly
    When starting "exo run" in my application directory
    Then my machine is running the services:
      | NAME             |
      | external-service |
    And it prints "all services online" in the terminal
