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
      | FIELD       | INPUT                           |
      | Description | serves HTML UI for the test app |
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
    And my application contains the file "web/service.yml" with the content:
    """
    name: web
    description: serves HTML UI for the test app

    setup: npm install --loglevel error --depth 0
    startup:
      command: node app
      online-text: web server is running

    messages:
      sends:
      receives:
    """


  Scenario: installing dependencies
    When running "exo setup" in this application's directory
    Then it has created the folders:
      | SERVICE | FOLDER       |
      | web     | node_modules |


  # Scenario: starting the application
  #   When starting "exo run" in this application's directory
  #   And waiting until I see "application ready" in the terminal
  #   Then requesting "http://localhost:3000" shows:
  #     """
  #     Welcome!
  #     """
  #   And I kill the server


  Scenario: adding the notes service
    When starting "exo add service todo exoservice-es6-mongodb" in this application's directory
    And entering into the wizard:
      | FIELD       | INPUT                   |
      | Description | stores the todo entries |
    And waiting until the process ends
    Then my application contains the file "todo/service.yml" with the content:
      """
      name: todo
      description: stores the todo entries

      setup: npm install --loglevel error --depth 0
      startup:
        command: node node_modules/exoservice/bin/exo-js
        online-text: online at port
      tests: node_modules/cucumber/bin/cucumber.js

      messages:
        receives:
          - todo.create
          - todo.create_many
          - todo.delete
          - todo.list
          - todo.read
          - todo.update
        sends:
          - todo.created
          - todo.created_many
          - todo.deleted
          - todo.listing
          - todo.details
          - todo.updated
      """
    When running "exo setup" in this application's directory
    And running "exo test" in this application's directory
    Then it prints "todo service works" in the terminal
    And it prints "web service has no tests, skipping" in the terminal
    And it prints "All tests passed" in the terminal


  Scenario: wiring up the web server to the todo service
    Given the file "web/app/controllers/index-controller.js":
      """
      class HomeController {

        constructor({send}) {
          this.send = send
        }

        index(req, res) {
          this.send('todo.list', {}, (todos) => {
            res.render('index', {todos})
          })
        }

      }


      module.exports = HomeController
      """
    And the file "web/app/views/index.jade":
      """
      extends layout

      block content

        h2 Exosphere Todos list
        p Your todos:
        ul
          each todo in todos
            li= todo

      h3 add a todo
      form(action="/todos" method="post")
        label text
        input(name="text")
        input(type="submit" value="add todo")
      """
    And the file "web/app/controllers/todos-controller.js":
      """
      class TodosController {

        constructor({send}) {
          this.send = send
        }

        create(req, res) {
          this.send('todo.create', { content: req.body.content }, () => {
            res.redirect('/')
          })
        }

      }
      module.exports = TodosController
      """
    And the file "web/app/routes.js":
      """
      module.exports = ({GET, resources}) => {
        GET('/', { to: 'home#index' })
        resources('todos', { only: ['create', 'destroy'] })
      }
      """
    And the file "web/service.yml":
      """
      name: web
      description: serves HTML UI for the test app

      setup: npm install --loglevel error --depth 0
      startup:
        command: node app
        online-text: web server is running

      messages:
        sends:
          - todo.list
        receives:
          - todo.listing
      """
    When starting "exo run" in this application's directory
    And waiting until I see "application ready" in the terminal
    Then requesting "http://localhost:3000" shows:
      """
      Exosphere Todos list
      Your todos:
      """
