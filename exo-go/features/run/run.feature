Feature: running Exosphere applications

  As an Exosphere developer
  I want to have an easy way to run an application on my developer machine
  So that I can test my app locally.

  Rules:
  - run "exo run" in the directory of your application to run it
  - this command boots up all the services and dependencies of the application


  # Scenario: booting a functioning Exosphere application
  #   Given I am in the root directory of the "running" example application
  #   When starting "exo run" in my application directory
  #   Then my machine has acquired the Docker images:
  #     | tmp_users        |
  #     | tmp_web          |
  #     | originate/exocom |
  #   Then the docker images have the following folders:
  #     | IMAGE             | FOLDER       |
  #     | tmp_mongo-service | node_modules |
  #     | tmp_web-server    | node_modules |
  #   Then it prints "all services online" in the terminal
  #   And my machine is running the services:
  #     | NAME  |
  #     | web   |


  Scenario: booting an Exosphere application with external docker images
    Given I am in the root directory of the "app-with-external-docker-images" example application
    When starting "exo run" in my application directory
    Then my machine has acquired the Docker images:
      | originate/test-web-server |
    And it prints "all services online" in the terminal
    And my machine is running the services:
      | NAME             |
      | external-service |
