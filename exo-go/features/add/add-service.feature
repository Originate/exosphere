Feature: interactive scaffolding

  As a tech lead
  I want to define the templates my team uses for creating new services
  So that they don't copy-and-past old code around.

  Rules:
  - templates for new services are located in the ".exosphere" folder
  - each subdirectory in that folder is a template
  - the templates are applied using boilr


  Scenario: adding a new service
    Given I am in the root directory of an empty application called "test app"
    And my application contains the template folder ".exosphere/foo" with the files:
      | NAME                     | CONTENT                      |
      | project.json             | { "Name": "foo" }            |
      | template/{{Name}}/foo.md | This is the {{Name}} service |
    When starting "exo add" in this application's directory
    And entering into the wizard:
      | FIELD                     | INPUT          |
      | Please choose a template: | 1              |
      | Protection Level:         | 1              |
    And waiting until the process ends
    Then my application now contains the file "application.yml" with the content:
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
    And my application now contains the file "test-service/service.yml" containing the text:
      """
      type: web-service
      description: testing
      author: tester
      """
    And my application now contains the file "test-service/foo.md" containing the text:
      """
      This is the test-service service
      """
