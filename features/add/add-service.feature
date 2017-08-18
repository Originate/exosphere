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
    And my application is a Git repository
    And my application has the templates:
      | NAME          | URL                                            |
      | exoservice-js | https://github.com/Originate/exoservice-js.git |
    When starting "exo add" in my application directory
    And entering into the wizard:
      | FIELD                         | INPUT          |
      | template                      | 1              |
      | serviceRole                   | ping-service   |
      | serviceType                   | ping-service   |
      | description                   | testing        |
      | author                        | tester         |
      | Protection Level              | 1              |
    And waiting until the process ends
    Then my application now contains the file "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0
      dependencies:
      - name: exocom
        version: 0.24.0
      services:
        public:
          ping-service:
            location: ./ping-service
        private: {}
      """
    And my application now contains the file "ping-service/service.yml" containing the text:
      """
      type: ping-service
      description: testing
      author: tester

      startup:
        command: node node_modules/exoservice/bin/exo-js
        online-text: online at port
      tests: node_modules/cucumber/bin/cucumber.js

      messages:
        receives:
          - ping
        sends:
          - pong

      dependencies:
      """
    And my application now contains the file "ping-service/README.md" containing the text:
      """
      # ping-service
      > testing
      """
