Feature: automatically restart a service when a change occurs

  As an Exosphere developer
  I want the services inside my application to automatically restart when I make a change
  So that I don't have to manually stop and restart the application

  Rules:
  - "exo run" will restart a service when a change occurs inside the service folder

  Scenario: automatically restart a service when a change occurs
    Given I am in the root directory of the "running" example application
    And starting "exo run" in my application directory
    And it prints "all dependencies online" in the terminal
    And it prints "all services online" in the terminal
    And it prints "watching ./mongo-service for changes" in the terminal
    When adding a file to "mongo-service" service folder
    Then the "users" service restarts
    And I stop all running processes
