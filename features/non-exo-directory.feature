Feature: running in a non-Exosphere directory

  As an application developer accidentally running "exo add" service in a directory that doesn't contain an Exosphere application
  I want to be given a friendly error message explaining that I have to run this command within an Exosphere application's directory
  So that I recognize my mistake and know how to run the command properly.

  Rules:

  - when running "exo add service" in a directory that doesn't contain a file "application.yml", the command prints: "Error: application.yml not found. Please run this command in the root directory of an Exosphere application."


  Scenario: add a service in a directory that is not an Exosphere application
    Given I am in an empty folder
    When starting "exo-add service" in this application's directory
    Then waiting until I see "Error: application.yml not found. Please run this command in the root directory of an Exosphere application." in the terminal
