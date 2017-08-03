Feature: Example test for scaffolded exo-services

  As a developer having scaffolded a new exo-service
  I want to see a running end-to-end test
  So that I can start TDDing my service without having to set up the test framework myself.

  - a new Exoservice contains an example Cucumber test for the "ping" command
  - run "exo test" in your app dir to run all tests for all services of your application
  - run "cucumber-js" in a service's directory to test only that service


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: uptime checking
    When receiving the "ping" command
    Then this service replies with a "pong" message
