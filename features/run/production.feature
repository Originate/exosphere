Feature: running Exosphere applications

  As an Exosphere developer
  I want to have an easy way to run an application in production mode
  So that I can test production mode before deploying

  Rules:
  - run "exo run --production" in the directory of your application to run it in production mode

  Scenario: booting an exosphere application in development mode
    Given I am in the root directory of the "static-asset-service" example application
    When starting "exo run" in my application directory
    And it prints "all services online" in the terminal
    Then http://localhost:8080 displays:
      """
      Application running in development mode
      """
    When starting "exo run --production" in my application directory
    And it prints "all services online" in the terminal
    Then http://localhost:8080 displays:
      """
      Application running in production mode
      """
