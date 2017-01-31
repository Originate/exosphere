Feature: scaffolding an ExpressJS html server written in ES6

  As a developer adding features to an Exosphere application
  I want to have an easy way to scaffold an empty service
  So that I don't waste time copy-and-pasting a bunch of code.

  - run "exo add service [name] htmlserver-express-es6" to add a new internal service to the current application
  - run "exo add service" to add the service interactively


  Scenario: calling with all command line arguments given
    Given I am in the root directory of an empty application called "test app"
    When running "exo-add service web html-server test-author htmlserver-express-es6 html description" in this application's directory
    Then my application contains the file "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0

      bus:
        type: exocom
        version: 0.16.1

      services:
        public:
          web:
            location: ./html-server
      """
    And my application contains the file "html-server/service.yml" with the content:
      """
      type: html-server
      description: description
      author: test-author

      # defines the commands to make the service runnable:
      # install its dependencies, compile it, etc.
      setup: yarn install

      # defines how to boot up the service
      startup:

        # the command to boot up the service
        command: node app

        # the string to look for in the terminal output
        # to determine when the service is fully started
        online-text: HTML server is running

      # the messages that this service will send and receive
      messages:
        sends:
        receives:

      # other services this service needs to run,
      # e.g. databases
      dependencies:

      docker:
        publish:
      """
    And my application contains the file "html-server/README.md" containing the text:
      """
      # TEST APP HTML Server
      > description
      """


  Scenario: calling with some command line arguments given
    Given I am in the root directory of an empty application called "test app"
    When starting "exo-add service web html-server test-author htmlserver-express-es6" in this application's directory
    And entering into the wizard:
      | FIELD                  | INPUT                           |
      | Description            | serves HTML UI for the test app |
      | Name of the data model |                                 |
    And waiting until the process ends
    Then my application contains the file "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0

      bus:
        type: exocom
        version: 0.16.1

      services:
        public:
          web:
            location: ./html-server
      """
    And my application contains the file "html-server/service.yml" with the content:
      """
      type: html-server
      description: serves HTML UI for the test app
      author: test-author

      # defines the commands to make the service runnable:
      # install its dependencies, compile it, etc.
      setup: yarn install

      # defines how to boot up the service
      startup:

        # the command to boot up the service
        command: node app

        # the string to look for in the terminal output
        # to determine when the service is fully started
        online-text: HTML server is running

      # the messages that this service will send and receive
      messages:
        sends:
        receives:

      # other services this service needs to run,
      # e.g. databases
      dependencies:

      docker:
        publish:
      """
    And my application contains the file "html-server/README.md" containing the text:
      """
      # TEST APP HTML Server
      > serves HTML UI for the test app
      """
