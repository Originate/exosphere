Feature: exo generate

  As a developer using different Exosphere commands
  I want to be able to easily generate needed docker-compose and terraform files
  So that I can keep these dependency files up-to-date

  Rules:
   - 'exo generate --check' should throw an error if files are missing or out of date


  Scenario: generate docker-compose and terraform files
    Given I am in the root directory of the "simple" example application
    When starting "exo generate docker-compose" in my application directory
    And waiting until the process ends
    Then my workspace contains the files:
      | docker-compose/run_development.yml |
      | docker-compose/run_production.yml  |
      | docker-compose/test.yml            |


  Scenario: generate docker-compose and terraform files
    Given I am in the root directory of the "simple" example application
    When starting "exo generate terraform" in my application directory
    And waiting until the process ends
    Then my workspace contains the files:
      | terraform/main.tf |


  Scenario: throwing an error when docker-compose files don't exist
    Given I am in the root directory of the "generate-check-dne" example application
    When starting "exo generate docker-compose --check" in my application directory
    Then it prints "'docker-compose/test.yml' does not exist. Please run 'exo generate'" in the terminal
    And it exits with code 1


  Scenario: throwing an error when docker-compose files are out of date
    Given I am in the root directory of the "generate-check-out-of-date-yml" example application
    When starting "exo generate docker-compose --check" in my application directory
    Then it prints "'docker-compose/test.yml' is out of date. Please run 'exo generate'" in the terminal
    And it exits with code 1


  Scenario: checking that docker-compose and terraform files exists and are up-to-date
    Given I am in the root directory of the "generate-check-good" example application
    When starting "exo generate docker-compose --check" in my application directory
    Then it exits with code 0
