Feature: testing an Exosphere application

  As a developer working on my Exosphere application
  I want to be able to run the tests for all my services in the application
  So that I know that everything is working


  Scenario: passing tests
    Given I am in the root directory of the "tests-passing" example application
    When starting "exo test" in my application directory
    Then I eventually see the following snippets:
      | Testing service 'users-service'  |
      | Testing service 'tweets-service' |
      | All tests passed                 |
    And it exits with code 0


  Scenario: failing tests
    Given I am in the root directory of the "tests-failing" example application
    When starting "exo test" in my application directory
    Then I eventually see the following snippets:
      | Testing service 'users-service'  |
      | Testing service 'tweets-service' |
      | 1 tests failed                   |
    And it exits with code 1
