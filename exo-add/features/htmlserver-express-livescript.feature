Feature: scaffolding an ExpressJS HTML service written in LiveScript

  As a developer adding features to an Exosphere application
  I want to have an easy way to scaffold an empty service
  So that I don't waste time copy-and-pasting a bunch of code.

  - run "exo add service [name] htmlserver-express-livescript" to add a new internal service to the current application
  - run "exo add service" to add the service interactively


  Scenario: scaffolding a LiveScript HTML server
    Given I am in the root directory of an empty application called "test app"
    When running "exo-add service --service-role=web --service-type=html-server --author=test-author --template-name=htmlserver-express-livescript --model-name=html --description=description  --protection-level=public" in this application's directory
    Then my application contains the file "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0

      dependencies:
        - type: exocom
          version: 0.22.0

      services:
        public:
          web:
            location: ./html-server
        private:
      """
    And my application contains the file "html-server/service.yml" with the content:
      """
      type: html-server
      description: description
      author: test-author

      setup: yarn install
      startup:
        command: ./node_modules/livescript/bin/lsc ./index.ls
        online-text: HTML server is running

      messages:
        sends:
        receives:

      dependencies:

      docker:
        ports:
          - '3000:3000'
      """
    And my application contains the file "html-server/README.md" containing the text:
      """
      # TEST APP HTML Server
      > description
      """
    And my application contains the file "html-server/.dockerignore" containing the text:
      """
      node_modules
      """
