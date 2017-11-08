Feature: running Exosphere applications that crash during setup

  Rules:
  - "exo run" exits on build errors


  Scenario: one of the set-up scripts for a service crashes
    Given I am in the root directory of the "failing-setup" example application
    When starting "exo run" in my application directory
    Then it prints "Cannot locate specified Dockerfile: Dockerfile.dev" in the terminal
    And it exits
