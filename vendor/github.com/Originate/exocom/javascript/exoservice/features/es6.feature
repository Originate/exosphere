Feature: Running services written in ES6

  As a developer familiar with JavaScript
  I want to be able to write and run exo-services written in ES6
  So that I can use the language that I'm familiar with for services.


  Scenario: running a service written in ES6
    Given an ExoCom instance
    And an instance of the "es6 service" service
    When receiving the "ping" message
    Then it acknowledges the received message
    And after a while it sends the "pong" message
