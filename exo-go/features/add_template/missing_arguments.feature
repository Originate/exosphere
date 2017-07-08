Feature: missing arguments

  Background:
    Given I am in the root directory of an empty application called "test app"
    And my application is a Git repository

  Scenario: attempting to run exo add-template without specifying template name or URL
    When starting "exo add-template" in my application directory
    Then I see:
      """
      not enough arguments
      """
    And it exits with code 1


  Scenario: attempting to run exo add-template without git URL
    When starting "exo add-template foo" in my application directory
    Then I see:
      """
      not enough arguments
      """
    And it exits with code 1