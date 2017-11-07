Feature: Mounting service directories in docker

  As an Exosphere developer
  I want the services inside my application to mount as volumes in the docker containers
  So that I can use existing tooling to restart services on change

  Rules:
  - exposed ports start at 3000 and icrement by 100 for each service

  Background:
    Given I am in the root directory of the "static-asset-service" example application

  Scenario: development
    When starting "exo run" in my application directory
    And it prints "nginx online" in the terminal
    Then http://localhost:3000 displays:
      """
      Application running in development mode
      """

  Scenario: production
    When starting "exo run --production" in my application directory
    And it prints "ngnix online" in the terminal
    Then http://localhost:3000 displays:
      """
      Application running in production mode
      """
