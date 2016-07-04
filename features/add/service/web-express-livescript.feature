Feature: scaffolding an ExpressJS web service written in LiveScript

  As a developer adding features to an Exosphere application
  I want to have an easy way to scaffold an empty service
  So that I don't waste time copy-and-pasting a bunch of code.

  - run "exo add service [name] web-express-livescript" to add a new internal service to the current application
  - run "exo add service" to add the service interactively


  Scenario: scaffolding a LiveScript web server
    Given I am in the root directory of an empty application called "test app"
    When running "exo add service web web-express-livescript description" in this application's directory
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
        command: ./node_modules/livescript/bin/lsc app
        online-text: web server is running

      messages:
        sends:
        receives:
      """
    And my application contains the file "web/README.md" containing the text:
      """
      # TEST APP Web Server
      > description
      """
