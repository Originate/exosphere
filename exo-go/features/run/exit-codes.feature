Feature: running Exosphere applications that crash during setup

  As an application developer
  I want to be able to recognize errors in my setup routine programmatically
  So that I can script application setup.

  Rules:
  - "exo run" returns the error code of the first failing subprocess


  Scenario: one of the set-up scripts for a service crashes
    Given I am in the root directory of the "failing-setup" example application
    When starting "exo run" in my application directory
    Then it exits with code 2
