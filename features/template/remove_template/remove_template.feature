Feature: removing existing templates

  When developing an exosphere application
	I want to be able to remove existing service templates from my application codebase
	So that I don't have to manually run the git command to remove the submodules

  - run "exo template remove <name>" to remove an existing service template from the application codebase


  Background:
    Given I am in the root directory of an empty application called "test-app"
    And my application is a Git repository


  Scenario: removing an existing service template
    Given my application has the templates:
      | NAME        | URL                                      |
      | boilr-spark | https://github.com/tmrts/boilr-spark.git |
    When running "exo template remove boilr-spark" in my application directory
    Then my git repository does not have any submodules


  Scenario: removing a non-existing service template
    When starting "exo template remove non-existing" in my application directory
    Then I eventually see:
      """
      Error: template does not exist
      """
    And it exits with code 1
