Feature: testing an Exosphere application

  As a developer working on my Exosphere application
  I want to be able to run the tests for all my services at once
  So that I know that everything is working


  Scenario: passing tests
    Given I am in the root directory of the "tests-passing" example application
    When starting "exo test" in my application directory
    Then it prints "running tests for users service" in the terminal
    And it prints "running tests for tweets service" in the terminal
    And it prints "users-service works" in the terminal
    And it prints "tweets-service works" in the terminal
    And it prints "All tests passed" in the terminal


  Scenario: failing tests
    Given I am in the root directory of the "tests-failing" example application
    When running "exo test"  in my application directory
    Then it prints "running tests for users service" in the terminal
    And it prints "running tests for tweets service" in the terminal
    And it prints "tweets-service works" in the terminal
    And it prints "users-service is broken" in the terminal
    And it prints "1 tests failed" in the terminal
