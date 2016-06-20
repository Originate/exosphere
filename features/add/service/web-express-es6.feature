Feature: scaffolding an ExpressJS web service written in ES6

  As a developer adding features to an Exosphere application
  I want to have an easy way to scaffold an empty service
  So that I don't waste time copy-and-pasting a bunch of code.

  - run "exo add service [name] web-express-es6" to add a new internal service to the current application
  - run "exo add service" to add the service interactively


  Scenario: calling without command-line arguments
    Given I am in the root directory of an empty application called "test app"
    When running "exo add service" in this application's directory
    And entering into the wizard:
      | FIELD                         | INPUT                           |
      | Name of the service to create | web                             |
      | Description                   | serves HTML UI for the test app |
      | Type                          |                                 |
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
    And my application contains the file "web/service.yml" with the content:
      """
      name: web
      description: serves HTML UI for the test app

      setup: npm install --loglevel error --depth 0
      startup:
        command: node app
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


  Scenario: calling with all command line arguments given
    Given I am in the root directory of an empty application called "test app"
    When running "exo add service web web-express-es6 description" in this application's directory
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
    And my application contains the file "web/service.yml" with the content:
      """
      name: web
      description: description

      setup: npm install --loglevel error --depth 0
      startup:
        command: node app
        online-text: all systems go

      messages:
        sends:
        receives:
      """
    And my application contains the file "web/README.md" containing the text:
      """
      # TEST APP Web Server
      > description
      """


  Scenario: calling with some command line arguments given
    Given I am in the root directory of an empty application called "test app"
    When running "exo add service web web-express-es6" in this application's directory
    And entering into the wizard:
      | FIELD                         | INPUT                           |
      | Description                   | serves HTML UI for the test app |
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
    And my application contains the file "web/service.yml" with the content:
      """
      name: web
      description: serves HTML UI for the test app

      setup: npm install --loglevel error --depth 0
      startup:
        command: node app
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
