Feature: scaffolding an ExoService written in LiveScript, backed by MongoDB

  As a LiveScript developer adding features to an Exosphere application
  I want to have an easy way to scaffold an ExoService backed by MongoDB
  So that I don't waste time copy-and-pasting a bunch of code.

  - run "exo add service <name> exoservice-ls-mongodb"
    to add a new exo-service written in LiveScript to the current application
  - run "exo add service" to add the service interactively


  Scenario: calling with all command line arguments
    Given I am in the root directory of an empty application called "test app"
    When running "exo add service user exoservice-ls-mongodb testing" in this application's directory
    Then my application contains the file "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0

      services:
        user:
          location: ./user
      """
    And my application contains the file "user/service.yml" with the content:
      """
      name: user
      description: testing

      setup: npm install --loglevel error --depth 0
      startup:
        command: node node_modules/exoservice/bin/exo-js
        online-text: all systems go
      tests: node_modules/cucumber/bin/cucumber.js

      messages:
        receives:
          - user.create
          - user.create_many
          - user.delete
          - user.list
          - user.read
          - user.update
        sends:
          - user.created
          - user.created_many
          - user.deleted
          - user.listing
          - user.details
          - user.updated
      """
    And my application contains the file "user/src/server.ls"
    And my application contains the file "user/README.md" containing the text:
      """
      # USER service
      > testing
      """
    When running "exo setup" in this application's directory
    And running "exo test" in this application's directory
    Then I see "All tests passed"
