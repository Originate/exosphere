Feature: scaffolding a custom template

  As a developer not using ExoService
  I want to have an easy way to use my own service templates
  So that I don't waste time copy-and-pasting a bunch of code.


  Background:
    Given I am in the root directory of an empty application called "test app"
    And my application initially contains the file ".exosphere/templates/add-service/custom1/service.yml" with the content:
      """
      type: _____serviceType_____
      description: _____description_____
      author: _____author_____

      setup: custom setup
      startup:
        command: custom command
        online-text: custom text
      tests: custom test

      messages:
        receives:
          - ping
        sends:
          - pong

      dependencies:
      """
    And my application initially contains the file ".exosphere/templates/add-service/custom1/README.md" with the content:
      """
      # _____serviceType@uppercase_____ service
      > _____description_____
      """
    And my application initially contains the file ".exosphere/templates/add-service/custom1/server.go" with the content:
      """
      package main

      func main() {
      	// run
      }
      """

  Scenario: calling with all command line arguments
    When running "exo-add service --service-role=users --service-type=users --author=test-author --template-name=custom1 --model-name=user --description=testing --protection-level=public" in this application's directory
    Then my application contains the file "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0

      dependencies:
        - type: exocom
          version: 0.22.1

      services:
        public:
          users:
            location: ./users
        private:
      """
    And my application contains the file "users/service.yml" with the content:
      """
      type: users
      description: testing
      author: test-author

      setup: custom setup
      startup:
        command: custom command
        online-text: custom text
      tests: custom test

      messages:
        receives:
          - ping
        sends:
          - pong

      dependencies:
      """
    And my application contains the file "users/README.md" containing the text:
      """
      # USERS service
      > testing
      """
    And my application contains the file "users/server.go" with the content:
      """
      package main

      func main() {
      	// run
      }
      """
