Feature: testing an Exosphere service

  As an Exosphere developer
  I want to be able to run my tests without mounting
  So that I can run tests in places where mounting isn't possible

  Rules:
  - when running "exo test", there is no mounting

  Scenario: testing a service
    Given I am in the root directory of the "fail-if-mounted" example application
    When starting "exo test" in the "test-service" directory
    Then it exits with code 0
