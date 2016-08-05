Feature: interactive scaffolding

  As a developer unsure about what the command-line arguments for the Exosphere scaffolder are
  I want to just call it as is and be asked for all relevant information
  So that I can levelage the scaffolder even if I'm not up to speed on it.


  Scenario: calling without command-line arguments
    Given I am in the root directory of an empty application called "test app"
    When starting "exo add service" in this application's directory
    And entering into the wizard:
      | FIELD                         | INPUT   |
      | Name of the service to create | web     |
      | Description                   | testing |
      | Type                          |         |
      | Name of the data model        | web     |
    And waiting until the process ends
    Then my application contains the file "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0

      services:
        web:
          location: ./web
      """
    And my application contains the file "web/service.yml" containing the text:
      """
      name: web
      description: testing
      """
    And my application contains the file "web/README.md" containing the text:
      """
      > testing
      """
      