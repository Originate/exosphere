Feature: interrupting Exosphere application tests

  As a developer test my Exosphere application
  I want to be able to interrupt the tests
  So that I don't have to waste time waiting for them to finish

  Rules:
  - press ctrl-c to interrupt running tests
  - Exosphere cleans up running test containers


  Scenario: clean up test containers
    Given I am in the root directory of the "tests-long-running" example application
    When starting "exo test" in my application directory
    Then it prints "Testing service 'test-service'" in the terminal
    And it prints "tests starting" in the terminal
    And I send an interrupt signal
    Then it prints "Stopping test-service" in the terminal
