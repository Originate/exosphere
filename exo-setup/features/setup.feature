Feature: Setup of Exosphere applications

  As a developer starting to work on an existing Exosphere application
  I want to be able to set up the application to run
  So that I can start developing without having to research how to set up each individual service.

  Rules:
  - run "exo setup" to set up all the services in the current application


  Scenario: set up the "test" application
    Given a freshly checked out "test" application
    When running "exo-setup" in this application's directory
    Then it finishes with exit code 0
    And it has created the folders:
      | SERVICE       | FOLDER       |
      | dashboard     | node_modules |
      | mongo-service | node_modules |
      | web-server    | node_modules |
    And it has created the files:
      | dashboard/Dockerfile     |
      | mongo-service/Dockerfile |
      | web-server/Dockerfile    |
    And it has generated the file "docker-compose.yml" with the content:
      """
      version: '3'
      services:
        exocom:
          image: 'originate/exocom:0.21.8'
          command: bin/exocom
          container_name: exocom
          environment:
            ROLE: exocom
            PORT: $EXOCOM_PORT
            SERVICE_ROUTES: >-
              [{role:web,receives:[users.listed,users.created],sends:[users.list,users.create]},{role:users,receives:[mongo.list,mongo.create],sends:[mongo.listed,mongo.created],namespace:mongo},{role:dashboard,receives:[users.listed,users.created],sends:[users.list]}]
        web:
          build: ./web-server
          container_name: web
          command: node_modules/livescript/bin/lsc server.ls
          environment:
            ROLE: web
            EXOCOM_HOST: exocom
            EXOCOM_PORT: $EXOCOM_PORT
          depends_on:
            - exocom
        users:
          build: ./mongo-service
          container_name: users
          command: node_modules/exoservice/bin/exo-js
          environment:
            ROLE: users
            EXOCOM_HOST: exocom
            EXOCOM_PORT: $EXOCOM_PORT
          depends_on:
            - exocom
        dashboard:
          build: ./dashboard
          container_name: dashboard
          command: node_modules/exoservice/bin/exo-js
          environment:
            ROLE: dashboard
            EXOCOM_HOST: exocom
            EXOCOM_PORT: $EXOCOM_PORT
          depends_on:
            - exocom
      """
    And it has acquired the Docker images:
      | dashboard      |
      | test_users     |
      | test_web       |
      | test_exocom    |

  Scenario: set up an application with external Docker images
    Given a freshly checked out "app-with-external-docker-images" application
    When running "exo-setup" in this application's directory
    Then it finishes with exit code 0
    And it has generated the file "docker-compose.yml" with the content:
      """
      version: '3'
      services:
        exocom:
          image: 'originate/exocom:0.21.8'
          command: bin/exocom
          environment:
            ROLE: exocom
            PORT: $EXOCOM_PORT
            SERVICE_ROUTES: '[]'
        external-service:
          image: originate/test-web-server
          depends_on:
            - exocom
      """
    And it has acquired the Docker images:
      | test-web-server |
