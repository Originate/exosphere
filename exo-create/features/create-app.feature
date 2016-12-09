Feature: scaffolding applications

  As a developer starting a new Exosphere application
  I want to have an easy way to scaffold an empty application shell
  So that I don't waste time looking up how Exosphere applications are structured and copy-and-pasting a bunch of code.

  - run "exo create application" to create a new application shell


  Scenario: creating an application
    When starting "exo-create application" in the terminal
    And entering into the wizard:
      | FIELD                             | INPUT              |
      | Name of the application to create | foo                |
      | Description                       | A test application |
      | Initial version                   | 0.0.0              |
    And waiting until I see "done" in the terminal
    Then my workspace contains the file "foo/application.yml" with content:
      """
      name: foo
      description: A test application
      version: 0.0.0

      services:
        public:
      """

  Scenario: creating an application with a name
    When starting "exo-create application foo" in the terminal
    And entering into the wizard:
      | FIELD                             | INPUT              |
      | Description                       | A test application |
      | Initial version                   |                    |
    And waiting until I see "done" in the terminal
    Then my workspace contains the file "foo/application.yml" with content:
      """
      name: foo
      description: A test application
      version: 0.0.1

      services:
        public:
      """

  Scenario: creating an application with a name and version number
    When starting "exo-create application foo 0.1" in the terminal
    And entering into the wizard:
      | FIELD                             | INPUT              |
      | Description                       | A test application |
    And waiting until I see "done" in the terminal
    Then my workspace contains the file "foo/application.yml" with content:
      """
      name: foo
      description: A test application
      version: 0.1

      services:
        public:
      """

  Scenario: creating an application with a name, version number and description
    When starting "exo-create application foo 0.1 A test application" in the terminal
    And waiting until I see "done" in the terminal
    Then my workspace contains the file "foo/application.yml" with content:
      """
      name: foo
      description: A test application
      version: 0.1

      services:
        public:
      """
