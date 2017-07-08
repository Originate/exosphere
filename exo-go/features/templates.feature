Feature: scaffolding applications

  When developing a fleet of similar applications
	I want to be able to use service templates from a remote location
	So that I don't have to copy-and-paste templates into all my application codebases.

  - run "exo fetch-templates" to fetch remote service templates and make those templates submodules
    of the application codebase


  Scenario: adding a new service template
  	Given I am in the root directory of an empty git application repository called "test app"
    When running "exo add-template boilr-spark https://github.com/tmrts/boilr-spark.git" in my application directory
    Then my application contains the directory ".exosphere/boilr-spark"
    And my git repository has a submodule ".exosphere/boilr-spark" with remote "https://github.com/tmrts/boilr-spark.git"
    And my application now contains the file "application.yml" with the content:
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
      templates:
        boilr-spark: https://github.com/tmrts/boilr-spark.git
      """


  Scenario: adding an existing service template
    Given I am in the root directory of an empty git application repository called "test app"
    When running "exo add-template boilr-spark https://github.com/tmrts/boilr-spark.git" in my application directory
    And starting "exo add-template boilr-spark https://github.com/tmrts/boilr-electron.git" in my application directory
    And entering into the wizard:
      | FIELD                         | INPUT          |
      | (y or n)                      | n              |
    Then I eventually see:
      """
      "boilr-spark" not updated
      """
    And my application now contains the file "application.yml" with the content:
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
      templates:
        boilr-spark: https://github.com/tmrts/boilr-spark.git
      """


  Scenario: updating an existing service template
    Given I am in the root directory of an empty git application repository called "test app"
    When running "exo add-template boilr-spark https://github.com/tmrts/boilr-spark.git" in my application directory
    And starting "exo add-template boilr-spark https://github.com/tmrts/boilr-electron.git" in my application directory
    And entering into the wizard:
      | FIELD                         | INPUT          |
      | (y or n)                      | y              |
    And waiting until the process ends
    Then my git repository has a submodule ".exosphere/boilr-spark" with remote "https://github.com/tmrts/boilr-electron.git"
    And my application now contains the file "application.yml" with the content:
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
      templates:
        boilr-spark: https://github.com/tmrts/boilr-electron.git
      """

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

  Scenario: removing an non-existing service template
    Given I am in the root directory of an empty git application repository called "test app"
    When starting "exo remove-template boilr-spark" in my application directory
    Then I eventually see:
      """
      Error: template does not exist
      """
    And it exits with code 1
