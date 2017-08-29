Feature: running Exosphere applications

  As an Exosphere developer
  I want to have an easy way to run an application on my developer machine
  So that I can test my app locally.

  Rules:
  - run "exo run" in the directory of your application to run it
  - this command boots up all the services and dependencies of the application

  Scenario: booting a functioning Exosphere application
    Given I am in the root directory of the "running" example application
    When starting "exo run" in my application directory
    Then it prints "setup complete" in the terminal
    And my machine has acquired the Docker images:
      | exampleapplication_users        |
      | exampleapplication_web          |
      | originate/exocom                |
    And the docker images have the following folders:
      | IMAGE                      | FOLDER       |
      | exampleapplication_users   | node_modules |
      | exampleapplication_web     | node_modules |
    And it prints "all dependencies online" in the terminal
    And it prints "all services online" in the terminal
    And my machine is running the services:
      | NAME  |
      | web   |
    And the network "exampleapplication_default" contains the running services:
      | NAME             |
      | web              |
      | users            |
      | originate/exocom |



  Scenario: booting an Exosphere application with external docker images
    Given I am in the root directory of the "app-with-external-docker-images" example application
    When starting "exo run" in my application directory
    Then it prints "setup complete" in the terminal
    And my machine has acquired the Docker images:
      | originate/test-web-server |
    And it prints "all dependencies online" in the terminal
    And it prints "all services online" in the terminal
    And my machine is running the services:
      | NAME             |
      | external-service |
    And the network "exampleapplication_default" contains the running services:
      | NAME             |
      | external-service |
      | originate/exocom |
