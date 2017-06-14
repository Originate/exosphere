Feature: running Exosphere applications

  As an Exosphere developer
  I want to have an easy way to run an application on my developer machine
  So that I can test my app locally.

  Rules:
  - run "exo run" in the directory of your application to run it
  - this command boots up all the services of the application


  Scenario: booting a functioning Exosphere application
    Given a running "simple" application
    Then my machine is running the services:
      | NAME |
      | web  |


  Scenario: booting a complex Exosphere application
    Given a running "running" application
    Then my machine is running ExoCom
    And my machine is running the services:
      | NAME  |
      | web   |
      | users |
    And ExoCom uses this routing:
      | ROLE  | SENDS                       | RECEIVES                    | NAMESPACE |
      | web   | users.list, users.create    | users.listed, users.created |           |
      | users | mongo.listed, mongo.created | mongo.list, mongo.create    | mongo     |
    When the web service broadcasts a "users.list" message
    Then the "mongo" service receives a "mongo.list" message
    And the "web" service receives a "users.listed" message


  Scenario: booting an Exosphere application with external dependencies
    Given a running "external-dependency" application
    Then my machine is running the services:
      | NAME  |
      | mongo |


  Scenario: Editing services of an Exosphere application
    Given a running "running" application
    Then my machine is running ExoCom
    And my machine is running the services:
      | NAME  |
      | web   |
      | users |
    When adding a file to the "users" service
    Then the "users" service restarts
    Then it prints "'users' restarted successfully" in the terminal


  Scenario: a service crashes during startup
    When trying to run the "crashing-service" application
    Then it prints "crasher exited with code 1" in the terminal


  Scenario: hiding logs from certain services or dependencies
    Given a running "silenced-running" application
    Then my machine is running the services:
      | NAME         |
      | web          |
      | users        |
      | exocom0.22.1 |
    And it prints "all services online" in the terminal
    And it prints "'web' is running" in the terminal
    And it does not print "'users' is running" in the terminal
    And it does not print "'exocom' is running" in the terminal
