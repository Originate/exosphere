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


  Scenario: verify that exo commands can be run by running "exo version"
    When starting "exo version" in the terminal
    Then I see "Exosphere version" in the terminal


  Scenario: setting up the application
    Given I am in an empty folder
    When starting "exo create application" in the terminal
    And entering into the wizard:
      | FIELD                             | INPUT              |
      | Name of the application to create | todo-app           |
      | Description                       | A todo application |
      | Initial version                   |                    |
      | ExoCom version                    | 0.16.3             |
    And waiting until the process ends
    Then my workspace contains the file "todo-app/application.yml" with content:
      """
      name: todo-app
      description: A todo application
      version: 0.0.1

      bus:
        type: exocom
        version: 0.16.3

      services:
        public:
      """


  Scenario: adding the html service
    Given I cd into "todo-app"
    When starting "exo add service --template-name=htmlserver-express-es6" in this application's directory
    And entering into the wizard:
      | FIELD                         | INPUT                           |
      | Role of the service to create | html-server                     |
      | Type of the service to create | html-server                     |
      | Description                   | serves HTML UI for the test app |
      | Author                        | test-author                     |
      | Name of the data model        |                                 |
    And waiting until the process ends
    Then my application contains the file "application.yml" with the content:
      """
      name: todo-app
      description: A todo application
      version: 0.0.1

      bus:
        type: exocom
        version: 0.16.3

      services:
        public:
          html-server:
            location: ./html-server
      """
    And my application contains the file "html-server/service.yml" with the content:
    """
    type: html-server
    description: serves HTML UI for the test app
    author: test-author

    # defines the commands to make the service runnable:
    # install its dependencies, compile it, etc.
    setup: yarn install

    # defines how to boot up the service
    startup:

      # the command to boot up the service
      command: node app

      # the string to look for in the terminal output
      # to determine when the service is fully started
      online-text: HTML server is running

    # the messages that this service will send and receive
    messages:
      sends:
      receives:

    # other services this service needs to run,
    # e.g. databases
    dependencies:

    docker:
      publish:
    """
    When running "exo setup" in this application's directory
    Then it has created the folders:
      | SERVICE     | FOLDER       |
      | html-server | node_modules |


  # Scenario: starting the application
  #   When starting "exo run" in this application's directory
  #   And waiting until I see "application ready" in the terminal
  #   Then requesting "http://localhost:3000" shows:
  #     """
  #     Welcome!
  #     """
  #   And I kill the server


  Scenario: adding the todo service
    When starting "exo add service --template-name=exoservice-es6-mongodb" in this application's directory
    And entering into the wizard:
      | FIELD                         | INPUT                   |
      | Role of the service to create | todo-service            |
      | Type of the service to create | todo-service            |
      | Description                   | stores the todo entries |
      | Author                        | test-author             |
      | Name of the data model        | todo                    |
    And waiting until the process ends
    Then my application contains the file "todo-service/service.yml" with the content:
      """
      type: todo-service
      description: stores the todo entries
      author: test-author

      setup: yarn install
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

      dependencies:
        mongo:
          docker_flags:
            volume: '-v {{EXO_DATA_PATH}}:/data/db'
            online_text: 'waiting for connections'
            port: '-p 27017:27017'
      """
    When running "exo setup" in this application's directory
    And running "exo test" in this application's directory
    Then it prints "todo-service works" in the terminal
    And it prints "html-server has no tests, skipping" in the terminal
    And it prints "All tests passed" in the terminal


  Scenario: wiring up the html server to the todo service
    Given the file "html-server/app/controllers/index-controller.js":
      """
      class IndexController {

        constructor({send}) {
          this.send = send
        }

        index(req, res) {
          this.send('todo.list', {}, (todos) => {
            res.render('index', {todos})
          })
        }

      }

      module.exports = IndexController
      """
    And the file "html-server/app/views/index.jade":
      """
      extends layout

      block content

        h2 Exosphere Todos list
        p Your todos:
        ul
          each todo in todos
            li= todo.text

        h3 add a todo
        form(action="/todos" method="post")
          label text
          input(name="text")
          input(type="submit" value="add todo")
      """
    And the file "html-server/app/controllers/todos-controller.js":
      """
      class TodosController {

        constructor({send}) {
          this.send = send
        }

        create(req, res) {
          this.send('todo.create', req.body, () => {
            res.redirect('/')
          })
        }

      }
      module.exports = TodosController
      """
    And the file "html-server/app/routes.js":
      """
      module.exports = ({GET, resources}) => {
        GET('/', { to: 'index#index' })
        resources('todos', { only: ['create', 'destroy'] })
      }
      """
    And the file "html-server/service.yml":
      """
      type: html-server
      description: serves HTML UI for the test app
      author: test-author

      setup: yarn install
      startup:
        command: node app
        online-text: HTML server is running

      messages:
        sends:
          - todo.create
          - todo.list
        receives:
          - todo.created
          - todo.listing

      dependencies:

      docker:
        publish:
          - '3000:3000'
      """
    When running "exo setup" in this application's directory
    And starting "exo run" in this application's directory
    And waiting until I see "all services online" in the terminal
    Then http://localhost:3000 displays:
      """
      Exosphere Todos list
      Your todos:
      """
    When adding a todo entry called "hello world" via the web application
    Then http://localhost:3000 displays:
      """
      Exosphere Todos list
      Your todos:
      hello world
      """
    And I stop all running processes
