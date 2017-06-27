Feature: displaying the version

  As an Exosphere user
  I want to have an easy and standardized way to determine the version of Exosphere I have installed on my machine
  So that I know whether I am up to date, and can work with customer support on investigating issues.

  - run "exo-go version" to see the version of Exosphere you are running


  Scenario: displaying the version
    When running "exo-go version"
    Then it prints:
      """
      Exosphere-Go v\d+\.\d+(\.\d+)?
      """
