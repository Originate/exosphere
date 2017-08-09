Feature: hiding logs from certain services or dependencies

  As an Exosphere developer
  I want to have an easy way to hide logs from specific services or dependencies
  So that I don't end up seeing irrelevant logs.

  Rules:
  - add "slient: true" to the service or dependency key in application.yml
  - "exo run" will not output any log from that service or dependency


  Scenario: hiding logs from certain services or dependencies
    Given I am in the root directory of the "silenced-running" example application
    When starting "exo run" in my application directory
    Then it prints "setup complete" in the terminal
    And it prints "all dependencies online" in the terminal
    And it prints "'web' is running" in the terminal
    And it prints "all services online" in the terminal
    And it does not print "'users' is running" in the terminal
    And it does not print "'exocom' is running" in the terminal
    And my machine is running the services:
      | NAME         |
      | web          |
      | users        |
      | exocom0.22.1 |
