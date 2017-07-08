Feature: fetching updates for existing templates

  When developing an exosphere application
	I want to be able to fetch the latest updates for the existing service templates in application codebase
	So that I don't have to manually run the git command to update the submodules

  - run "exo fetch-templates" to add a remote service template to the application codebase

  Scenario: fetching updates for all existing service templates
    Given I am in the root directory of an empty git application repository called "test app"
    When running "exo add-template boilr-spark https://github.com/tmrts/boilr-spark.git" in my application directory
    And running "exo fetch-templates" in my application directory
    Then it prints "done" in the terminal
