Feature: missing arguments

  Background:
    Given I am in the root directory of an empty application called "test-app"
    And my application is a Git repository

  Scenario: attempting to run `exo template fetch` without specifying template name
    When starting "exo template fetch" in my application directory
    Then I see:
      """
      wrong number of arguments
      """
    And it exits with code 1
