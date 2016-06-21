Feature: testing an Exosphere application

  As a developer working on my Exosphere application
  I want to be able to run the tests for all my services at once
  So that I know that everything is working


  Scenario: passing tests
    Given a set-up "tests-passing" application
    When running "exo test" in this application's directory
    Then I see "running tests for users service"
    And I see "running tests for tweets service"
    And I see "All tests passed"


  Scenario: failing tests
    Given a set-up "tests-failing" application
    When running "exo test" in this application's directory
    Then I see "running tests for users service"
    And I see "running tests for tweets service"
    And I see "Tests failed"
