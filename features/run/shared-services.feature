Feature: Shared service directories

  As an Exosphere developer
  I want to be able to reuse certain files in every service
  So that I don't have to duplicate code

  Rules:
  - Shared directories are listed under `shared-directories` in `application.yml`
  - Each shared directory is copied into every other service directory


  Scenario: shared directories are copied into service directories
    Given I am in the root directory of the "complex-setup-app" example application
    And starting "exo run" in my application directory
    And it prints "Creating complexsetupapp" in the terminal
    Then every service contains a shared directory "shared"
