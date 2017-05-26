Feature: interactive scaffolding

  As a developer unsure about what the command-line arguments for the Exosphere scaffolder are
  I want to just call it as is and be asked for all relevant information
  So that I can levelage the scaffolder even if I'm not up to speed on it.


  Scenario: calling without command-line arguments
    Given I am in the root directory of an empty application called "test app"
    When starting "exo-add service" in this application's directory
    And entering into the wizard:
      | FIELD                         | INPUT   |
      | Role of the service to create | web     |
      | Type of the service to create | web     |
      | Description                   | testing |
      | Author                        | tester  |
      | Template                      |         |
      | Name of the data model        | web     |
      | Protection level              | public  |
    And waiting until the process ends
    Then my application contains the file "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0

      bus:
        type: exocom
        version: 0.21.7

      services:
        public:
          web:
            location: ./web
        private:
      """
    And my application contains the file "web/service.yml" containing the text:
      """
      type: web
      description: testing
      author: tester
      """
    And my application contains the file "web/README.md" containing the text:
      """
      > testing
      """

