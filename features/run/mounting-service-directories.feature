Feature: Mounting service directories in docker

  As an Exosphere developer
  I want the services inside my application to mount as volumes in the docker containers
  So that I can use existing tooling to restart services on change

  Rules:
  - "exo run" mounts service directories in /mnt on the docker container

  Scenario: changes made in the service directory appear in docker
    Given I am in the root directory of the "frontend-with-webpack" example application
    And starting "exo run" in my application directory
    And it prints "all services online" in the terminal
    Then http://localhost:8080 displays:
      """
      Hello world
      """
    When modifying frontend-service/src/index.html to "Foobar"
    And modifying frontend-service/src/index.js to "console.log()"
    Then it prints "webpack: Compiled successfully" in the terminal
    And http://localhost:8080 displays:
      """
      Foobar
      """

  Scenario: run with `--no-mount` to disable mounting
    Given I am in the root directory of the "fail-if-mounted" example application
    When starting "exo run --no-mount" in my application directory
    Then it prints "test-service exited with code 0" in the terminal
