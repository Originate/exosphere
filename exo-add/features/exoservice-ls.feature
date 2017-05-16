Feature: scaffolding an ExoService written in LiveScript

  As a LiveScript developer adding features to an Exosphere application
  I want to have an easy way to scaffold an empty ExoService
  So that I don't waste time copy-and-pasting a bunch of code.

  - run "exo add service <name> exoservice-ls"
    to add a new exo-service written in LiveScript to the current application
  - run "exo add service" to add the service interactively


  Scenario: calling with all command line arguments
    Given I am in the root directory of an empty application called "test app"
    When running "exo-add service --service-role=users --service-type=users --author=test-author --template-name=exoservice-ls --model=user --description=testing" in this application's directory
    Then my application contains the file "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0

      bus:
        type: exocom
        version: 0.21.7

      services:
        public:
          users:
            location: ./users
      """
    And my application contains the file "users/service.yml" with the content:
      """
      type: users
      description: testing
      author: test-author

      setup: yarn install
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
    And my application contains the file "users/src/server.ls" with the content:
      """
      module.exports =

        before-all: (done) ->
          # TODO: add asynchronous init code here, or delete the whole block
          done!


        # replies to the "ping" command
        ping: (_, {reply}) ->
          reply 'pong'
      """
    And my application contains the file "users/README.md" containing the text:
      """
      # USERS service
      > testing
      """
