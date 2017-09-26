Feature: displaying the version

  As an Exosphere user
  I want to have an easy and standardized way to determine the version of Exosphere I have installed on my machine
  So that I know whether I am up to date, and can work with customer support on investigating issues.

  - run "exo version" to see the version of Exosphere you are running


  Scenario: displaying the version
    When running "exo version" in the terminal
    Then it prints "Exosphere v0.26.2" in the terminal
