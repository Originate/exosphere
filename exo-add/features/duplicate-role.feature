Feature: Attempting to add a duplicate service

  As a LiveScript developer adding features to an Exosphere application
  I want Exosphere to alert me when attempting to create a duplicate service
  So that I don't create conflicting services in my application


  Scenario: Adding a service-role that already exists
    Given I am in the root directory of an empty application called "test app"
    When running "exo-add service users test-author exoservice-ls user testing" in this application's directory
    Then my application contains the file "application.yml" with the content:
      """
      name: test app
      description: Empty test application
      version: 1.0.0

      bus:
        type: exocom
        version: 0.16.1

      services:
        public:
          users:
            location: ./users
      """
    When running "exo-add service users test-author exoservice-ls user testing" in this application's directory
    Then it prints "Service 'users' already exists in this application" in the terminal

