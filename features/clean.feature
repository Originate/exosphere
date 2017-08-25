Feature: cleaning dangling Docker images

  As a developer using Exosphere
  I want to be able to easily clean up my Docker workspace
  So that unused images and volumes do not take up space on my Docker VM

  Rules:
  - run "exo clean" in the terminal in any directory to remove dangling
    Docker images and volumes on your machine
  - this command does not remove non-dangling Docker images/volumes


  Background:
    Given I am in the root directory of the "clean-containers" example application

  Scenario: cleaning a machine with both dangling and non-dangling Doker images
    Given my machine has both dangling and non-dangling Docker images and volumes
    When running "exo clean" in my application directory
    Then it prints "removed all dangling images" in the terminal
    And it prints "removed all dangling volumes" in the terminal
    And it has non-dangling images
    And it does not have dangling images
    And it does not have dangling volumes

  Scenario: cleaning a machine with running application and service test containers
    Given my machine has running application and service test containers
    And my machine has running third party containers
    # When running "exo clean" in the terminal
    # Then it prints "removed application containers" in the terminal
    # Then it prints "removed service test containers" in the terminal
    # And it does not stop any third party containers
    # And it stops and removes application and service test containers


  # Scenario: cleaning a machine with stopped application and service test containers
  #   Given my machine has stopped application adn service test containers
  #   And my machine has stopped third party containers
  #   When running "exo clean" in the terminal
  #   Then it prints "removed application containers" in the terminal
  #   Then it prints "removed service test containers" in the terminal
  #   And it does not remove any third party containers
  #   And it removes application and service test containers
