@e2e
Feature: Following the tutorial

  As a person learning Exosphere
  I want that the whole tutorial works end to end
  So that I can follow along with the examples without getting stuck on bugs.

  AC:
  - all steps in the tutorial work when executed one after the other

  Notes:
  - The steps only do quick verifications.
    Full verifications are in the individual specs for the respective step.
  - You can not run individual scenarios here,
    you always have to run the whole feature.


  Scenario: setting up the application
    Given I am in an empty folder
    When starting "exo create application" in the terminal
    And entering into the wizard:
      | FIELD                             | INPUT              |
      | Name of the application to create | todo-app           |
      | Description                       | A todo application |
      | Initial version                   |                    |
    And waiting until the process ends
    Then my workspace contains the file "todo-app/application.yml" with content:
      """
      name: todo-app
      description: A todo application
      version: 0.0.1

      services:
      """


  Scenario: adding the web service
    Given I cd into "todo-app"
    When starting "exo add service web web-express-es6" in this application's directory
    And entering into the wizard:
      | FIELD                         | INPUT                           |
      | Description                   | serves HTML UI for the test app |
    And waiting until the process ends
    Then my application contains the file "application.yml" with the content:
      """
      name: todo-app
      description: A todo application
      version: 0.0.1

      services:
        web:
          location: ./web
      """

  Scenario: installing dependencies
    When running "exo setup" in this application's directory
    Then it has created the folders:
      | SERVICE | FOLDER       |
      | web     | node_modules |


  Scenario: starting the application
    When starting "exo run" in this application's directory
    And waiting until I see "application ready"
    Then requesting "http://localhost:3000" shows:
      """
      Welcome!
      """
    And I kill the server
