Feature: running Exosphere applications that crash during startup

  As an Exosphere developer
  I want to be notified when any service of my application crashes during startup

  Rules:
  - "exo run" prints "[service-name] exited with exited code 1" when a service crashes


  Scenario: a service crashes during startup
    Given I am in the root directory of the "crashing-service" example application
    When starting "exo run" in my application directory
    Then it prints "crasher exited with code 1" in the terminal
    And I stop all running processes
