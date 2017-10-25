Feature: application dependencies

  As an Exosphere developer
  I want to have an easy way to run application wide dependencies
  So that all my services can use these dependencies.

  Rules:
  - list application dependencies in the "dependencies" key of "application.yml"
  - these dependencies get booted up before the application services


  Scenario: booting an application that uses Exocom
    Given I am in the root directory of the "simple" example application
    When starting "exo run" in my application directory
    Then it prints "setup complete" in the terminal
    And my machine has acquired the Docker images:
      | originate/exocom |
    And it prints "ExoCom online at port" in the terminal
    And my machine is running the services:
      | NAME         |
      | exocom0.26.1 |


  Scenario: booting an application that uses NATS
    Given I am in the root directory of the "nats" example application
    When starting "exo run" in my application directory
    Then it prints "setup complete" in the terminal
    And my machine has acquired the Docker images:
      | nats0.9.6 |
    And it prints "Listening for route connections" in the terminal
    And my machine is running the services:
      | NAME      |
      | nats0.9.6 |
