Feature: cleaning dangling Docker images

  As a developer using Exosphere
  I want to be able to easily clean up my Docker workspace
  So that unused images and volumes do not take up space on my Docker VM

  Rules:
  - run "exo clean" in the terminal in any directory to remove dangling
    Docker images and volumes on your machine
  - this command does not remove non-dangling Docker images/volumes


  Scenario: cleaning a machine with both dangling and non-dangling Doker images
    Given my machine has both dangling and non-dangling Docker images
    And it has the Docker images:
      | tmp_web |
      | <none>  |
    When running "exo-clean" in the terminal
    Then it prints "removed all dangling images" in the terminal
    And it prints "removed all dangling volumes" in the terminal
    And it has the Docker images:
      | tmp_web |
    And it does not have the Docker images:
      | <none>  |
