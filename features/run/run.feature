Feature: running Exosphere applications

  As an Exosphere developer
  I want to have an easy way to run an application on my developer machine
  So that I can test my app locally.

  Rules:
  - run "exo run" in the directory of your application to run it
  - this command boots up all the services and dependencies of the application

  Scenario: booting a functioning Exosphere application
    Given I am in the root directory of the "running" example application
    When starting "exo run" in my application directory
    Then it prints "online at port" in the terminal
    And it prints "web server running at port" in the terminal
    And my machine has acquired the Docker images:
      | originate/exocom:0.26.1 |
    And the docker images have the following folders:
      | IMAGE         | FOLDER       |
      | running_users | node_modules |
      | running_web   | node_modules |
    And my machine contains the network "running_default"
    And the network "running_default" contains the running services:
      | NAME                   |
      | running_web_1          |
      | running_users_1        |
      | running_exocom0.26.1_1 |


  Scenario: booting an Exosphere application with external docker images
    Given I am in the root directory of the "app-with-external-docker-images" example application
    When starting "exo run" in my application directory
    Then it prints "web server running at port" in the terminal
    And my machine has acquired the Docker images:
      | originate/test-web-server:0.0.1 |
    And my machine contains the network "appwithexternaldockerimages_default"
    And the network "appwithexternaldockerimages_default" contains the running services:
      | NAME                                           |
      | appwithexternaldockerimages_external-service_1 |
      | appwithexternaldockerimages_exocom0.26.1_1     |


  Scenario: booting a functioning Exosphere application from a service directory
    Given I am in the root directory of the "simple" example application
    When starting "exo run" in the "web" directory
    Then it prints "online at port" in the terminal
    And it prints "web server running at port" in the terminal
    And my machine has acquired the Docker images:
      | originate/exocom:0.26.1 |
    And my machine contains the network "simple_default"
    And the network "simple_default" contains the running services:
      | NAME                  |
      | simple_web_1          |
      | simple_exocom0.26.1_1 |
