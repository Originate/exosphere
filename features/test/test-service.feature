Feature: testing an Exosphere service

  As an Exosphere developer working with a variety of services written in different languages
  I want to have a standardized way to run all tests for a service
  So that I don't have to waste time investigating how to use language-specific tools.

  Rules:
  - when running "exo test" in a service directory, it runs the test for the respective service

  Scenario: testing a service
    Given I am in the root directory of the "tests-passing" example application
    When starting "exo test" in the "tweets-service" directory
    Then it prints "Testing service 'tweets-service'" in the terminal
    And it exits with code 0
    And it does not print "Testing service 'users-service'" in the terminal

  Scenario: executing "exo test" in a directory that isn't a service or application
    Given I am in the root directory of a non-exosphere application
    When starting "exo test" in my application directory
    Then it doesn't run any tests
    And it exits with code 1
