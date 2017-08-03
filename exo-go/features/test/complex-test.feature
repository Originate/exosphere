Feature: testing an Exosphere application

  As a developer working on my Exosphere application
  I want to be able to run the tests that requires
  So that I know that everything is working


  Scenario: tests that require extra setup
    Given I am in the root directory of the "tests-cucumber" example application
    When starting "exo test" in my application directory
    Then it prints "Testing service 'exoservice'" in the terminal
    And I eventually see:
      """
      1 scenario (1 passed)
      """
    And it prints "'exoservice' tests passed" in the terminal
    And it prints "All tests passed" in the terminal
