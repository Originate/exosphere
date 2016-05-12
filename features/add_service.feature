Feature: scaffolding services

  As a developer adding features to an Exosphere application
  I want to have an easy way to scaffold an empty service
  So that I don't waste time copy-and-pasting a bunch of code.

  - run "exo add service" to add a new internal service to the current application


  Scenario: adding a web service
    Given I am in the root directory of an empty application called "test app"
    When running "exo add service" in this application's directory
    And entering into the wizard:
      | FIELD                         | INPUT                           |
      | Name of the service to create | web                             |
      | Description                   | serves HTML UI for the test app |
      | Type                          |                                 |
    And waiting until I see "done"
    Then my application contains the file "web/config.yml" with the content:
      """
      name: web
      description: serves HTML UI for the test app

      startup:
        command: node_modules/.bin/lsc app
        online-text: all systems go

      messages:
        sends:
        receives:
      """
    And my application contains the file "web/README.md" containing the text:
      """
      # TEST APP Web Server
      > serves HTML UI for the test app
      """
