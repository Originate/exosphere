Feature: testing an Exosphere application

  As a developer working on my Exosphere application
  I want to be able to run the tests for all my services at once
  So that I know that everything is working


  Scenario: passing tests
    Given a set-up "tests-passing" application
    When running "exo test" in this application's directory
    Then it prints "running tests for users service" in the terminal
    And it prints "running tests for tweets service" in the terminal
    And it prints "users-service works" in the terminal
    And it prints "tweets-service works" in the terminal
    And it prints "All tests passed" in the terminal


  Scenario: failing tests
    Given a set-up "tests-failing" application
    When running "exo test" in this application's directory
    Then it prints "running tests for users service" in the terminal
    And it prints "running tests for tweets service" in the terminal
    And it prints "tweets-service works" in the terminal
    And it prints "users-service is broken" in the terminal
    And it prints "Tests failed" in the terminal
