Feature: interactive scaffolding

  As a developer unsure about what the command-line arguments for the Exosphere scaffolder are
  I want to just call it as is and be asked for all relevant information
  So that I can levelage the scaffolder even if I'm not up to speed on it.


  Scenario: calling without command-line arguments
    Given I am in the root directory of an empty application called "test app"
    And my application contains the template folder ".exosphere/foo" with the files:
      | NAME                     | CONTENT                      |
      | project.json             | { "Name": "foo" }            |
      | template/{{Name}}/foo.md | This is the {{Name}} service |
    When starting "exo-add" in this application's directory
    And entering into the wizard:
      | FIELD                         | INPUT          |
      | template                      | 1              |
      | Name                          | test-service   |
      | ServiceType                   | web-service    |
      | Description                   | testing        |
      | Author                        | tester         |
      | Protection Level              | 1              |
    And waiting until the process ends
    Then my application contains the file "application.yml" with the content:
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
        private: {}
      """
    And my application contains the file "test-service/service.yml" containing the text:
      """
      type: web-service
      description: testing
      author: tester
      """
    And my application contains the file "test-service/foo.md" containing the text:
      """
      This is the test-service service
      """
