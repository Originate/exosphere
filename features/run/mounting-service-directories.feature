Feature: Mounting service directories in docker

  As an Exosphere developer
  I want the services inside my application to mount as volumes in the docker containers
  So that I can use existing tooling to restart services on change

  Rules:
  - "exo run" mounts service directories in /mnt on the docker container

  Scenario: changes made in the service directory appear in docker
    Given I am in the root directory of the "frontend-with-webpack" example application
    And starting "exo run" in my application directory
    And it prints "webpack: Compiled successfully" in the terminal
    Then http://localhost:3000 displays:
      """
      Hello world
      """
    When modifying frontend-service/src/index.html to "Foobar"
    And modifying frontend-service/src/index.js to "console.log()"
    Then it prints "webpack: Compiled successfully" in the terminal
    And http://localhost:3000 displays:
      """
      Foobar
      """
