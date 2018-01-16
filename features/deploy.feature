Feature: deploy command

  Scenario Outline: 'exo deploy' with not enough arguments
    When starting "exo deploy <SUBCOMMAND>" in the terminal
    Then I see:
      """
      Error: Wrong number of arguments
      """
    And it exits with code 1

  Examples:
    | SUBCOMMAND     |
    | infrastructure |
    | services       |

  Scenario Outline: 'exo deploy' with too many arguments
    When starting "exo deploy <SUBCOMMAND> a b" in the terminal
    Then I see:
      """
      Error: Wrong number of arguments
      """
    And it exits with code 1

  Examples:
    | SUBCOMMAND     |
    | infrastructure |
    | services       |

  Scenario Outline: 'exo deploy' with invalid remote id
    Given I am in the root directory of the "simple" example application
    When starting "exo deploy <SUBCOMMAND> production" in the terminal
    Then I see:
      """
      Invalid remote environment id: production. Valid remote environment ids: qa
      """
    And it exits with code 1

  Examples:
    | SUBCOMMAND     |
    | infrastructure |
    | services       |
