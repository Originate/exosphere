Feature: scaffolding an ExoService written in ES6

  As a JavaScript developer adding features to an Exosphere application
  I want to have an easy way to scaffold an empty ExoService
  So that I don't waste time copy-and-pasting a bunch of code.

  - run "exo add service <name> exoservice-es6"
    to add a new exo-service written in ES6 to the current application
  - run "exo add service" to add the service interactively


  Scenario: calling with all command line arguments
    Given I am in the root directory of an empty application called "test app"
    When running "exo-add service users test-author exoservice-es6 user testing" in this application's directory
    Then my application contains the file "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0

      services:
        public:
          users:
            docker_image: test-author/users
            location: ./users
      """
    And my application contains the file "users/service.yml" with the content:
      """
      title: users
      description: testing
      author: test-author

      setup: npm install --loglevel error --depth 0
      startup:
        command: node node_modules/exoservice/bin/exo-js
        online-text: online at port
      tests: node_modules/cucumber/bin/cucumber.js

      messages:
        receives:
          - ping
        sends:
          - pong

      docker:
        link:
      """
    And my application contains the file "users/src/server.js" with the content:
      """
      module.exports = {

        beforeAll: (done) => {
          // TODO: add asynchronous init code here, or delete the whole block
          done()
        },

        // Replies to the "ping" command
        ping: (_, {reply}) => {
          reply('pong')
        }

      }
      """
    And my application contains the file "users/README.md" containing the text:
      """
      # USERS service
      > testing
      """
