Feature: adding remote templates

  When developing a fleet of similar applications
	I want to be able to use service templates from a remote location
	So that I don't have to copy-and-paste templates into all my application codebases.

  - run "exo template add <name> <git-url>" to add a remote service template to the application codebase


  Background:
    Given I am in the root directory of an empty application called "test app"
    And my application is a Git repository


  Scenario: adding a new service template
    When running "exo template add boilr-spark https://github.com/tmrts/boilr-spark.git" in my application directory
    Then my application contains the directory ".exosphere/boilr-spark"
    And my git repository has a submodule ".exosphere/boilr-spark" with remote "https://github.com/tmrts/boilr-spark.git"


  Scenario: the service template already exists and is not overwritten
    Given my application has the templates:
      | NAME        | URL                                      |
      | boilr-spark | https://github.com/tmrts/boilr-spark.git |
    When starting "exo template add boilr-spark https://github.com/tmrts/boilr-electron.git" in my application directory
    Then I see:
      """
      The template "boilr-spark" already exists
      """
    And it exits with code 1
