Feature: exo generate

  As a developer using different Exosphere commands
  I want to be able to easily generate needed docker-compose and terraform files
  So that I can keep these dependency files up-to-date


  Scenario: generate docker-compose and terraform files
    Given I am in the root directory of the "complex-setup" example application
    When starting "exo generate" in my application directory
    Then it prints "generating docker-compose and terraform files" in the terminal
    And my workspace contains the files:
      | docker-compose/run_development.yml |
      | docker-compose/run_production      |
      | docker-compose/test.yml            |
      | terraform/main.tf                  |

