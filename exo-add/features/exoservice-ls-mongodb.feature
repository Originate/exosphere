Feature: scaffolding an ExoService written in LiveScript, backed by MongoDB

  As a LiveScript developer adding features to an Exosphere application
  I want to have an easy way to scaffold an ExoService backed by MongoDB
  So that I don't waste time copy-and-pasting a bunch of code.

  - run "exo add service <name> exoservice-ls-mongodb"
    to add a new exo-service written in LiveScript to the current application
  - run "exo add service" to add the service interactively


  Scenario: calling with all command line arguments
    Given I am in the root directory of an empty application called "test app"
    When running "exo-add service --service-role=user-service --service-type=user-service --author=test-author --template-name=exoservice-ls-mongodb --model-name=user --description=testing  --protection-level=public" in this application's directory
    Then my application contains the file "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0

      dependencies:
        - type: exocom
          version: 0.21.7

      services:
        public:
          user-service:
            location: ./user-service
        private:
      """
    And my application contains the file "user-service/service.yml" with the content:
      """
      type: user-service
      description: testing
      author: test-author

      setup: yarn install
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
          - user.details
          - user.listing
          - user.updated

      dependencies:
        mongo:
          dev:
            image: 'mongo'
            version: '3.4.0'
            volumes:
              - '{{EXO_DATA_PATH}}:/data/db'
            ports:
              - '27017:27017'
      """
    And my application contains the file "user-service/src/server.ls"
    And my application contains the file "user-service/README.md" containing the text:
      """
      # USER-SERVICE
      > testing
      """
    And my application contains the file "user-service/.dockerignore" containing the text:
      """
      node_modules
      """
