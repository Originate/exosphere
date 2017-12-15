Feature: configure command

  Scenario: 'exo configure' with not enough arguments
    When starting "exo configure" in the terminal
    Then I see:
      """
      Error: Wrong number of arguments
      """
    And it exits with code 1

  Scenario: 'exo configure' with too many arguments
    When starting "exo configure a b" in the terminal
    Then I see:
      """
      Error: Wrong number of arguments
      """
    And it exits with code 1

  Scenario: 'exo configure' with invalid remote id
    Given I am in the root directory of the "simple" example application
    When starting "exo configure production" in the terminal
    Then I see:
      """
      Invalid remote environment id: production. Valid remote environment ids: qa
      """
    And it exits with code 1

  Scenario: 'exo configure create' with not enough arguments
    When starting "exo configure create" in the terminal
    Then I see:
      """
      Error: Wrong number of arguments
      """
    And it exits with code 1

  Scenario: 'exo configure create' with too many arguments
    When starting "exo configure create a b" in the terminal
    Then I see:
      """
      Error: Wrong number of arguments
      """
    And it exits with code 1

  Scenario: 'exo configure create' with invalid remote id
    Given I am in the root directory of the "simple" example application
    When starting "exo configure create production" in the terminal
    Then I see:
      """
      Invalid remote environment id: production. Valid remote environment ids: qa
      """
    And it exits with code 1

  Scenario: 'exo configure delete' with not enough arguments
    When starting "exo configure delete" in the terminal
    Then I see:
      """
      Error: Wrong number of arguments
      """
    And it exits with code 1

  Scenario: 'exo configure delete' with too many arguments
    When starting "exo configure delete a b" in the terminal
    Then I see:
      """
      Error: Wrong number of arguments
      """
    And it exits with code 1

  Scenario: 'exo configure delete' with invalid remote id
    Given I am in the root directory of the "simple" example application
    When starting "exo configure delete production" in the terminal
    Then I see:
      """
      Invalid remote environment id: production. Valid remote environment ids: qa
      """
    And it exits with code 1

  Scenario: 'exo configure read' with not enough arguments
    When starting "exo configure read" in the terminal
    Then I see:
      """
      Error: Wrong number of arguments
      """
    And it exits with code 1

  Scenario: 'exo configure read' with too many arguments
    When starting "exo configure read a b" in the terminal
    Then I see:
      """
      Error: Wrong number of arguments
      """
    And it exits with code 1

  Scenario: 'exo configure read' with invalid remote id
    Given I am in the root directory of the "simple" example application
    When starting "exo configure read production" in the terminal
    Then I see:
      """
      Invalid remote environment id: production. Valid remote environment ids: qa
      """
    And it exits with code 1

  Scenario: 'exo configure update' with not enough arguments
    When starting "exo configure update" in the terminal
    Then I see:
      """
      Error: Wrong number of arguments
      """
    And it exits with code 1

  Scenario: 'exo configure update' with too many arguments
    When starting "exo configure update a b" in the terminal
    Then I see:
      """
      Error: Wrong number of arguments
      """
    And it exits with code 1

  Scenario: 'exo configure update' with invalid remote id
    Given I am in the root directory of the "simple" example application
    When starting "exo configure update production" in the terminal
    Then I see:
      """
      Invalid remote environment id: production. Valid remote environment ids: qa
      """
    And it exits with code 1
