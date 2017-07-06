Feature: attempting to add a duplicate service

  As a developer adding features to an Exosphere application
  I want Exosphere to alert me when attempting to create a duplicate service
  So that I don't create conflicting services in my application


  Scenario: adding a service-role that already exists
    Given I am in the directory of "test app" application containing a "test-service" service
    When starting "exo add" in my application directory
    And entering into the wizard:
      | FIELD                         | INPUT          |
      | template                      | 1              |
      | Name                          | test-service   |
      | ServiceType                   | web-service    |
      | Description                   | testing        |
      | Author                        | tester         |
      | Protection Level              | 1              |
    Then I eventually see:
      """
      Service test-service already exists in this application
      """
    And it exits with code 1
