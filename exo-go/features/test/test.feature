Feature: testing an Exosphere application

  As a developer working on my Exosphere application
  I want to be able to run the tests for all my services in the application
  So that I know that everything is working


  Scenario: passing tests
    Given I am in the root directory of the "tests-passing" example application
    When starting "exo test" in my application directory
    Then it prints "Testing service 'users-service'" in the terminal
    And it prints "Testing service 'tweets-service'" in the terminal
    And it prints "'users-service' tests passed" in the terminal
    And it prints "'tweets-service' tests passed" in the terminal
    And it prints "All tests passed" in the terminal


  Scenario: failing tests
    Given I am in the root directory of the "tests-failing" example application
    When starting "exo test" in my application directory
    Then it prints "Testing service 'users-service'" in the terminal
    And it prints "Testing service 'tweets-service'" in the terminal
    And it prints "'tweets-service' tests passed" in the terminal
    And it prints "'users-service' tests failed" in the terminal
    And it prints "1 tests failed" in the terminal
