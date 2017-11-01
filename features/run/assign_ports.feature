Feature: Mounting service directories in docker

  As an Exosphere developer
  I want the services inside my application to mount as volumes in the docker containers
  So that I can use existing tooling to restart services on change

  Rules:
  - exposed ports start at 3000 and icrement by 100 for each service

  Scenario: development
    Given I am in the root directory of the "frontend-with-webpack" example application
    And starting "exo run" in my application directory
    And it prints "webpack: Compiled successfully" in the terminal
    Then http://localhost:3000 displays:
      """
      Hello world
      """

  Scenario: production
    Given I am in the root directory of the "frontend-with-webpack" example application
    And starting "exo run --production" in my application directory
    And it prints "attaching to ngnix online" in the terminal
    Then http://localhost:3000 displays:
      """
      Hello world
      """
