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


  Scenario: tutorial
    ########################################
    # Printing the exosphere version
    ########################################
    When running "exo version" in my application directory
    Then it prints "Exosphere v0.25.1" in the terminal

    ########################################
    # Setting up the application
    ########################################
    Given I am in an empty folder
    When starting "exo create" in my application directory
    And entering into the wizard:
      | FIELD              | INPUT              |
      | AppName            | todo-app           |
      | AppDescription     | A todo application |
      | AppVersion         |                    |
      | ExocomVersion      | 0.26.1             |
    And waiting until the process ends
    Then my workspace contains the file "todo-app/application.yml" with content:
      """
      name: todo-app
      description: A todo application
      version: 0.0.1

      development:
        dependencies:
          - name: exocom
            version: 0.26.1

      services:
        public:
        private:
      """
    And my workspace contains the empty directory "todo-app/.exosphere"
    And I cd into "todo-app"
    And running "git init" in my application directory

    ########################################
    # Adding the html service
    ########################################
    Given I add the "exosphere-htmlserver-express" template
    When starting "exo add" in my application directory
    And entering into the wizard:
      | FIELD                         | INPUT                            |
      | template                      | 1                                |
      | serviceRole                   | html-server                      |
      | appName                       | test-app                         |
      | serviceType                   | html-server                      |
      | description                   | serves HTML UI for the test app  |
      | author                        | test-author                      |
      | Protection Level              | 1                                |
    And waiting until the process ends
    Then my application now contains the file "application.yml" with the content:
      """
      name: todo-app
      description: A todo application
      version: 0.0.1
      development:
        dependencies:
        - name: exocom
          version: 0.26.1
      services:
        public:
          html-server:
            location: ./html-server
        private: {}
        worker: {}
      """
    And my application now contains the file "html-server/service.yml" with the content:
    """
    type: html-server
    description: serves HTML UI for the test app
    author: test-author

    startup:
      online-text: HTML server is running

    development:
      scripts:
        run: node ./index.js
    """
    When starting "exo run" in my application directory
    And waiting until I see "setup complete" in the terminal
    Then the docker images have the following folders:
      | IMAGE               | FOLDER       |
      | todoapp_html-server | node_modules |
    And I stop all running processes

    ########################################
    # adding the todo service
    ########################################
    Given I add the "exoservice-js-mongodb" template
    When starting "exo add" in my application directory
    And entering into the wizard:
      | FIELD                         | INPUT                    |
      | template                      | 1                        |
      | serviceRole                   | todo-service             |
      | description                   | stores the todo entries  |
      | serviceType                   | todo-service             |
      | author                        | test-author              |
      | modelName                     | todo                     |
      | EXO_DATA_PATH                 |                          |
      | Protection Level              | 1                        |
    And waiting until the process ends
    Then my application now contains the file "todo-service/service.yml" with the content:
      """
      type: todo-service
      description: stores the todo entries
      author: test-author

      startup:
        online-text: online at port

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

      development:
        scripts:
          run: node src/server.js
          test: node_modules/cucumber/bin/cucumber.js
        dependencies:
          - name: 'mongo'
            version: '3.4.0'
            config:
              volumes:
                - '{{EXO_DATA_PATH}}:/data/db'
              ports:
                - '27017:27017'
              online-text: 'waiting for connections'
      """

    ########################################
    # wiring up the html server to the todo service
    ########################################
    Given the file "html-server/app/controllers/index-controller.js":
      """
      class IndexController {

        constructor({send}) {
          this.send = send
        }

        index(req, res) {
          this.send('todo.list', {}, (messageName, payload) => {
            if (messageName === 'todo.listing') {
              res.render('index', {todos: payload})
            } else {
              res.sendStatus(500)
            }
          })
        }

      }

      module.exports = IndexController
      """
    And the file "html-server/app/views/index.pug":
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

      startup:
        online-text: HTML server is running

      messages:
        sends:
          - todo.create
          - todo.list
        receives:
          - todo.created
          - todo.listing

      docker:
        ports:
          - '3000:3000'

      development:
        scripts:
          run: node ./index.js
      """
    When starting "exo run" in my application directory
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

      * hello world
      """
    And I stop all running processes
