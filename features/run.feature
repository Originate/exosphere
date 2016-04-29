Feature: running Exosphere applications

  As an Exosphere developer
  I want to have an easy way to run an application on my developer machine
  So that I can test my app locally.

  Rules:
  - run "exo run" in the directory of your application to run it
  - this command boots up all dependencies


  Scenario: running the "test" application
    Given a set-up "test" application
    When running "exo run" in this application's directory
    And waiting until I see "application ready"
    Then my machine is running ExoCom
    And ExoCom uses this routing:
      | COMMAND       | SENDERS        | RECEIVERS      |
      | users.list    | web, dashboard | users          |
      | users.listed  | users          | web, dashboard |
      | users.create  | web            | users          |
      | users.created | users          | web, dashboard |
    And my machine is running the services:
      | NAME      |
      | web       |
      | users     |
      | dashboard |
