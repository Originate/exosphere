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
    And it does not print "The following tests failed:" in the terminal


  Scenario: failing tests
    Given I am in the root directory of the "tests-failing" example application
    When starting "exo test" in my application directory
    Then I eventually see the following snippets:
      | Testing service 'users-service'  |
      | Testing service 'tweets-service' |
      | The following tests failed:      |
    And it exits with code 1
