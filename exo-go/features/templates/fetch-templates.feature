Feature: scaffolding applications

  When developing a fleet of similar applications
	I want to be able to use service templates from a remote location
	So that I don't have to copy-and-paste templates into all my application codebases.

  - run "exo fetch-templates" to fetch remote service templates and make those templates submodules
    of the application codebase


  Scenario: fetching and updating remote service templates
  	Given I am in the root directory of an empty git application repository called "test app"
  	And my application contains the file "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0
      dependencies:
      - name: exocom
        version: 0.22.1
      services:
        public:
          test-service:
            location: ./test-service
        private:
      templates:
        boilr-spark: https://github.com/tmrts/boilr-spark.git
      """
    When starting "exo fetch-templates" in my application directory
    And waiting until the process ends
    Then my application contains the directory ".exosphere/boilr-spark"
    And my git repository has ".exosphere/boilr-spark" as a submodule
    And the git URL of ".exosphere/boilr-spark" is "https://github.com/tmrts/boilr-spark.git"
    When I update "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0
      dependencies:
      - name: exocom
        version: 0.22.1
      services:
        public:
          test-service:
            location: ./test-service
        private:
      templates:
        boilr-spark: https://github.com/tmrts/boilr-electron.git
      """
    And starting "exo fetch-templates" in my application directory
    And waiting until the process ends
    Then my git repository has ".exosphere/boilr-spark" as a submodule
    And the git URL of ".exosphere/boilr-spark" is "https://github.com/tmrts/boilr-electron.git"
