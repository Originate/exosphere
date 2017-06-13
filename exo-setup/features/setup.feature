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

  Scenario: set up the application with silenced services and dependencies
    Given a freshly checked out "silenced-running" application
    When running "exo-setup" in this application's directory
    Then it does not print "users  setup finished"
