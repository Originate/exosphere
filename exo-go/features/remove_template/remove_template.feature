Feature: removing existing templates

  When developing an exosphere application
	I want to be able to remove existing service templates from my application codebase
	So that I don't have to manually run the git command to remove the submodules

  - run "exo remove-template" to remove an existing service template from the application codebase


  Scenario: removing an existing service template
    Given I am in the root directory of an empty git application repository called "test app"
    When running "exo add-template boilr-spark https://github.com/tmrts/boilr-spark.git" in my application directory
    And running "exo remove-template boilr-spark" in my application directory
    Then my application now contains the file "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0
      dependencies:
      - name: exocom
        version: 0.22.1
      services:
        public: {}
        private: {}
      """


  Scenario: removing a non-existing service template
    Given I am in the root directory of an empty git application repository called "test app"
    When starting "exo remove-template boilr-spark" in my application directory
    Then I eventually see:
      """
      Error: template does not exist
      """
    And it exits with code 1
