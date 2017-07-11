Feature: fetching updates for existing templates

  When developing an exosphere application
	I want to be able to fetch the latest updates for the existing service templates in application codebase
	So that I don't have to manually run the git command to update the submodules

  - run "exo template fetch" to fetch the latest updates for the existing service templates

  Scenario: fetching updates for all existing service templates
    Given I am in the root directory of an empty application called "test app"
    And my application is a Git repository
    And my application has the templates:
      | NAME        | URL                                      |
      | boilr-spark | https://github.com/tmrts/boilr-spark.git |
    When running "exo template fetch" in my application directory
    Then it prints "done" in the terminal
