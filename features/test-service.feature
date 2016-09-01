Feature: testing an Exosphere service

  As an Exosphere developer working with a variety of services written in different languages
  I want to have a standardized way to run all tests for a service
  So that I don't have to waste time investigating how to use language-specific tools.

  Rules:
  - when running "exo test" in a service directory, it runs the test for the respective service

  @verbose
  Scenario: testing a service
    Given I am in the "tests-passing/tweets-service" directory
    When running "exo-test"
    Then it only runs tests for "tweets-service"

  Scenario: executing "exo test" in a directory that isn't a service or application
    Given I am in the "tests-passing/tweets-service/src" created directory
    When running "exo-test"
    Then it doesn't run any tests
