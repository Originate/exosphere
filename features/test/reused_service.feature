Feature: testing an Exosphere service

  As an Exosphere developer working with a variety of services written in different languages
  I want to have a standardized way to run all tests for a service
  So that I don't have to waste time investigating how to use language-specific tools.

  Rules:
  - when running "exo test" in a service directory, it runs the test for the respective service

  Scenario: testing a service
    Given I am in the root directory of the "reused_service" example application
    When starting "exo test" in the "test-service" directory
    Then it prints "Testing service 'test-service'" in the terminal
    And it prints "'test-service' tests passed" in the terminal
    And it exits with code 0

  Scenario: testing all services
    Given I am in the root directory of the "reused_service" example application
    When starting "exo test" in my application directory
    Then it prints "Testing service 'test-service'" in the terminal
    And it prints "'test-service' tests passed" in the terminal
    And it exits with code 0
