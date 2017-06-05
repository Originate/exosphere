Feature: cleaning dangling Docker images

  As an Exosphere developer
  I want to have an easy way to clean up dangling <none>:<none> Docker
  images on my developer machine created by an Exosphere application
  So that I don't end up with a lot of irrelevant Docker images after 
  I build my application multiple times.

  Rules:
  - run "exo clean" in the terminal in any directory to remove dangling
    <none>:<none> Docker images and volumes on your machine
  - this command does not remove non-dangling Docker images/volumes


  Scenario: cleaning a machine with both dangling and non-dangling Doker images
    Given a running "running" application
    When adding a file to the "users" service
    When setting up "running" application again
    Then my machine has a number of dangling and non-dangling Docker images
    And it has the Docker images:
      | tmp_web |
      | tmp_users |
      | <none> |
    When running "exo-clean" in the terminal
    Then it prints "removed all dangling images" in the terminal
    Then it prints "removed all dangling volumes" in the terminal
    Then only dangling Docker images are removed
    And it has the Docker images:
      | tmp_web |
      | tmp_users |
    And it does not have the Docker images:
      | <none> |


  Scenario: cleaning a machine that may or may not have (dangling) Doker images
    Then my machine has a number of dangling and non-dangling Docker images
    When running "exo-clean" in the terminal
    Then it prints "removed all dangling images" in the terminal
    Then it prints "removed all dangling volumes" in the terminal
    Then only dangling Docker images are removed
    And it does not have the Docker images:
      | <none> |
