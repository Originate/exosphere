Feature: missing arguments

  Scenario: attempting to run exo remove-template without specifying template name
    Given I am in the root directory of an empty git application repository called "test app"
    When starting "exo remove-template" in my application directory
    Then I eventually see:
      """
      not enough arguments
      """
    And it exits with code 1
