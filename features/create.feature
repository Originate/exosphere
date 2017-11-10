Feature: scaffolding applications

  As a developer starting a new Exosphere application
  I want to have an easy way to scaffold an empty application shell
  So that I don't waste time looking up how Exosphere applications are structured and copy-and-pasting a bunch of code.

  - run "exo create" to create a new application shell


  Scenario: creating an application
    When starting "exo create" in the terminal
    And entering into the wizard:
      | FIELD                    | INPUT              |
      | AppName                  | foo                |
      | AppDescription           | A test application |
      | AppVersion               | 0.0.0              |
      | ExocomVersion            | latest             |
    Then I eventually see "done" in the terminal
    And my workspace contains the file "foo/application.yml" with content:
      """
      name: foo
      description: A test application
      version: 0.0.0

      development:
        dependencies:
          - name: exocom
            version: latest

      services:
        public:
        private:
        worker:
      """
    And my workspace contains the empty directory "foo/.exosphere/service_templates"
