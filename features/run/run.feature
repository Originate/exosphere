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
      | running_users    |
      | running_web      |
      | originate/exocom |
    And the docker images have the following folders:
      | IMAGE         | FOLDER       |
      | running_users | node_modules |
      | running_web   | node_modules |
    And it prints "all dependencies online" in the terminal
    And it prints "all services online" in the terminal
    And my machine is running the services:
      | NAME  |
      | web   |
    And my machine contains the network "running_default"
    And the network "running_default" contains the running services:
      | NAME         |
      | web          |
      | users        |
      | exocom0.24.0 |



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
    And my machine contains the network "exampleapplication_default"
    And the network "exampleapplication_default" contains the running services:
      | NAME             |
      | external-service |
      | exocom0.24.0     |
