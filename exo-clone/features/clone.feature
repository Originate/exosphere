Feature: cloning an Exosphere application and its services

  As an application developer
  I want to download the complete source code for an application onto my machine
  So that I can work on it without manually cloning multiple repositories.

  - "exo clone" downloads the Exosphere application at the given URL onto the user's machine
  - the URL must point to a Git repository containing an Exosphere application,
    i.e. have an application.yml file in the root of the repository
  - all services referenced by the application.yml file also get cloned onto the machine


  Scenario: cloning an application with external services
    Given I am in an empty folder
    And source control contains the services "web" and "users"
    And source control contains a repo "app" with a file "application.yml" and the content:
      """
      name: Example application
      description: Demonstrates basic Exosphere application startup
      version: '1.0'

      services:
        web:
          local: ./web
          origin: ../origins/web
        users:
          local: ./users
          origin: ../origins/users
      """
    When running "exo-clone origins/app" in the terminal
    Then it creates the files:
      | app/application.yml   |
      | app/web/service.yml   |
      | app/users/service.yml |

  Scenario: clone an application with services that reside above the application directory.
    Given I am in an empty folder
    And source control contains a "web" service
    And source control contains a repo "app" with a file "application.yml" and the content:
      """
      name: Example application
      description: Demonstrates Exosphere application with service residing above app directory
      version: '1.0'

      services:
        web:
          local: ../web
          origin: ../origins/web
      """
    When running "exo-clone origins/app" in the terminal
    Then it creates the files:
      | app/application.yml |
      | web/service.yml     |


  Scenario: Cloning an application with a nonexisting service.
    Given I am in an empty folder
    And source control contains a "web" service
    And source control contains a repo "app" with a file "application.yml" and the content:
      """
      name: Example application
      description: contains a non-existing service
      version: '1.0'

      services:
        nonexisting:
          local: ./nonexisting
          origin: ../origins/nonexisting
        web:
          local: ./web
          origin: ../origins/web
      """
    When running "exo-clone origins/app" in the terminal
    Then I get the error "fatal: repository '../origins/nonexisting' does not exist"
    And it prints "web  done" in the terminal
    And it prints "exo-clone  Failed" in the terminal
    And no new files or folders have been created

  Scenario: the application directory already exists
    Given I am in an empty folder
    And source control contains a "users" service
    And source control contains a repo "app" with a file "application.yml" and the content:
      """
      name: Example application
      description: is going to be cloned into an already existing directory
      version: '1.0'

      services:
        users:
          local: ./users
          origin: ../origins/users
      """
    And my workspace already contains the folder "app"
    When trying to run "exo-clone origins/app"
    Then I get the error "exo-clone  fatal: destination path 'app' already exists"
    And no new files or folders have been created


  Scenario: a service directory already exists
    Given I am in an empty folder
    And source control contains a "users" service
    And source control contains a repo "app" with a file "application.yml" and the content:
      """
      name: Example application
      description: contains a service that would be cloned into an already existing directory
      version: '1.0'

      services:
        users:
          local: ../users
          origin: ../origins/users
      """
    And my workspace already contains the folder "users"
    When running "exo-clone origins/app" in the terminal
    Then I get the error "users  Service cloning failed"
    And it prints "exo-clone  Some services failed to clone or were invalid Exosphere services." in the terminal
    And it prints "exo-clone  Failed" in the terminal
