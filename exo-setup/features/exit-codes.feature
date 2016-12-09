@verbose
Feature: setup script returns with error code of failing service

  As an application developer
  I want to be able to recognize errors in my setup routine programmatically
  So that I can script application setup.

  Rules:
  - exo setup returns the error code 0 if all setup work,
    otherwise the error code of the first failing subprocess


  Scenario: one of the set-up scripts for a service crashes
    Given a freshly checked out "failing-setup" application
    When running "exo-setup"
    Then it finishes with exit code 3

  Scenario: service setup runs smoothly
    Given a freshly checked out "test" application
    When running "exo-setup"
    Then it finishes with exit code 0
