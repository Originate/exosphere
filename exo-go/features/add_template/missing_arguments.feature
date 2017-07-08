Feature: missing arguments

  Scenario: attempting to run exo add-template without specifying template name or URL
    Given I am in the root directory of an empty git application repository called "test app"
    When starting "exo add-template" in my application directory
    Then I eventually see:
      """
      not enough arguments
      """
    And it exits with code 1


  Scenario: attempting to run exo add-template without git URL
    Given I am in the root directory of an empty git application repository called "test app"
    When starting "exo add-template foo" in my application directory
    Then I eventually see:
      """
      not enough arguments
      """
    And it exits with code 1