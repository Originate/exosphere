Feature: test templates

  When developing an exosphere service template
	I want to be able to test it to make sure it can be used by exocom
	So that I don't have broken templates

  - run "exo template test" in the directory of an exosphere template to test
    adding it to an application and run the tests


  Scenario: success
    Given I am in the root directory of the "good" example template
    When running "exo template test" in my template directory
    Then it prints "Template passes" in the terminal
    And it exits with code 0


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
