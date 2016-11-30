Feature: Linting Exosphere applications

  As an Exosphere developer
  I want to have an easy way to confirm all messages being sent are received somewhere
  So that I know that all my messages are handled.

  Rules:
  - run "exo lint" in the directory of your application to lint it
  - if a sent message is not received anywhere, it prints an error message
  - if a received message it not sent anywhere, it prints an error message

  Scenario: functioning application
    Given I am in the directory of an application with the services:
      | NAME  | SENDS                      | RECEIVES                   |
      | web   | user.create, user.delete   | user.created, user.deleted |
      | users | user.created, user.deleted | user.create, user.delete   |
    When running "exo-lint" in this application's directory
    Then it prints "Lint passed" in the terminal


  Scenario: some sent messages are not listened to
    Given I am in the directory of an application with the services:
      | NAME  | SENDS                      | RECEIVES                   |
      | web   | user.create, user.delete   | user.created, user.deleted |
      | users | user.created, user.deleted | user.create                |
    When running "exo-lint" in this application's directory
    Then it prints "The following messages are sent but not received:" in the terminal
    Then it prints "web  user.delete" in the terminal


  Scenario: some received messages are not sent
    Given I am in the directory of an application with the services:
      | NAME  | SENDS                      | RECEIVES                   |
      | web   | user.create, user.delete   | user.created, user.deleted |
      | users | user.deleted               | user.create, user.delete   |
    When running "exo-lint" in this application's directory
    Then it prints "The following messages are received but not sent:" in the terminal
    And it prints "web  user.created" in the terminal
