Feature: Setup of Exosphere applications

  As a developer starting to work on an existing Exosphere application
  I want to be able to set up the application to run
  So that I can start developing without having to research how to set up each individual service.

  Rules:
  - run "exo setup" to set up all the services in the current application


  Scenario: set up the "test" application
    Given a freshly checked out "test" application
    When running "exo-setup" in this application's directory
    Then it has created the folders:
      | SERVICE       | FOLDER       |
      | dashboard     | node_modules |
      | mongo-service | node_modules |
      | web-server    | node_modules |
    And it has created the files:
      | dashboard/Dockerfile     |
      | mongo-service/Dockerfile |
      | web-server/Dockerfile    |
    And it has acquired the Docker images:
      | tmp_dashboard    |
      | tmp_users        |
      | tmp_web          |
      | originate/exocom |
    And ExoCom uses this routing:
      | ROLE      | SENDS                       | RECEIVES                    | NAMESPACE |
      | web       | users.list, users.create    | users.listed, users.created |           |
      | users     | mongo.listed, mongo.created | mongo.list, mongo.create    | mongo     |
      | dashboard | users.list                  | users.listed, users.created |           |

  Scenario: set up an application with external Docker images
    Given a freshly checked out "app-with-external-docker-images" application
    When running "exo-setup" in this application's directory
    Then it has acquired the Docker images:
      | originate/test-web-server |
    And ExoCom uses this routing:
      | ROLE             | SENDS                       | RECEIVES                    | NAMESPACE |
      | external-service | users.list, users.create    | users.listed, users.created |           |

  Scenario: set up an application with services shared dependencies
    Given a freshly checked out "shared-dependency" application
    When running "exo-setup" in this application's directory
    Then my application contains the file "tmp/docker-compose.yml" with the content:
      """
      version: '3'
      services:
        exocom0.22.1:
          image: 'originate/exocom:0.22.1'
          command: bin/exocom
          container_name: exocom0.22.1
          environment:
            ROLE: exocom
            PORT: $EXOCOM_PORT
            SERVICE_ROUTES: >-
              [{"role":"web","receives":["users.listed","users.created"],"sends":["users.list","users.create"]},{"role":"html-server","receives":["todo.created"],"sends":["todo.create"]},{"role":"users","receives":["mongo.list","mongo.create"],"sends":["mongo.listed","mongo.created"],"namespace":"mongo"}]
        web:
          build: ../web-server
          container_name: web
          command: node_modules/.bin/lsc server.ls
          ports:
            - '4000:4000'
          links:
            - 'mongo3.4.0:mongo'
          environment:
            ROLE: web
            EXOCOM_HOST: exocom0.22.1
            EXOCOM_PORT: $EXOCOM_PORT
            MONGO: mongo
          depends_on:
            - mongo3.4.0
            - exocom0.22.1
        html-server:
          build: ../html-server
          container_name: html-server
          command: echo "does not run"
          ports:
            - '3000:3000'
          environment:
            ROLE: html-server
            EXOCOM_HOST: exocom0.22.1
            EXOCOM_PORT: $EXOCOM_PORT
          depends_on:
            - exocom0.22.1
        users:
          build: ../mongo-service
          container_name: users
          command: node_modules/exoservice/bin/exo-js
          links:
            - 'mongo3.4.0:mongo'
          environment:
            ROLE: users
            EXOCOM_HOST: exocom0.22.1
            EXOCOM_PORT: $EXOCOM_PORT
            MONGO: mongo
          depends_on:
            - mongo3.4.0
            - exocom0.22.1
        mongo3.4.0:
          image: 'mongo:3.4.0'
          container_name: mongo3.4.0
          ports:
            - '5000:5000'
          volumes:
            -
      """
