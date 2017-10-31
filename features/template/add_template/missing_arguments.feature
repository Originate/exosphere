Feature: missing arguments

  Background:
    Given I am in the root directory of an empty application called "test-app"
    And my application is a Git repository

  Scenario: attempting to run 'exo template add' without specifying template name or URL
    When starting "exo template add" in my application directory
    Then I see:
      """
      wrong number of arguments
      """
    And it exits with code 1


  Scenario: attempting to run 'exo template add' without git URL
    When starting "exo template add foo" in my application directory
    Then I see:
      """
      wrong number of arguments
      """
    And it exits with code 1
