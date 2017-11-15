Feature: exo generate

  As a developer using different Exosphere commands
  I want to be able to easily generate needed docker-compose and terraform files
  So that I can keep these dependency files up-to-date


  Scenario: generate docker-compose and terraform files
    Given I am in the root directory of the "simple" example application
    When starting "exo generate" in my application directory
    And waiting until the process ends
    Then my workspace contains the files:
      | docker-compose/run_development.yml |
      | docker-compose/run_production.yml  |
      | docker-compose/test.yml            |
      | terraform/main.tf                  |
