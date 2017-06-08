Feature: running Exosphere applications with dependencies

  As an Exosphere developer
  I want to have an easy way to run application wide dependencies
  So that all my services can use these dependencies.

  Rules:
  - run "exo run" in the directory of your application to run it
  - this command boots up all the dependencies of the application


  Scenario: booting an Exosphere application with Exocom as a listed dependency
    Given a running "simple" application
    Then my machine is running the services:
      | NAME         |
      | exocom0.22.0 |


  Scenario: booting an Exosphere application with Nats as a listed dependency
    Given a running "nats" application
    Then my machine is running the services:
      | NAME      |
      | nats0.9.6 |
