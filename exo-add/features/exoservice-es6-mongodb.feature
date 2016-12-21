Feature: scaffolding an ExoService written in ES6, backed by MongoDB

  As a JavaScript developer adding features to an Exosphere application
  I want to have an easy way to scaffold an ExoService backed by MongoDB
  So that I don't waste time copy-and-pasting a bunch of code.

  - run "exo add service <name> exoservice-es6-mongodb"
    to add a new exo-service written in ES6 to the current application
  - run "exo add service" to add the service interactively


  @verbose
  Scenario: calling with all command line arguments
    Given I am in the root directory of an empty application called "test app"
    When running "exo-add service user-service test-author exoservice-es6-mongodb user testing" in this application's directory
    Then my application contains the file "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0

      services:
        public:
          user-service:
            docker_image: test-author/user-service
            location: ./user-service
      """
    And my application contains the file "user-service/service.yml" with the content:
      """
      title: user-service
      description: testing
      author: test-author

      setup: npm install --loglevel error --depth 0
      startup:
        command: node node_modules/exoservice/bin/exo-js
        online-text: online at port
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

      dependencies:
        - 'mongo'
      """
    And my application contains the file "user-service/src/server.js"
    And my application contains the file "user-service/README.md" containing the text:
      """
      # USER-SERVICE
      > testing
      """
