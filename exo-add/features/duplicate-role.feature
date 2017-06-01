Feature: Attempting to add a duplicate service

  As a developer adding features to an Exosphere application
  I want Exosphere to alert me when attempting to create a duplicate service
  So that I don't create conflicting services in my application


  Scenario: Adding a service-role that already exists
    Given I am in the directory of an application containing a "users" service of type "users-service"
    When trying to run "exo-add service --service-role=users --service-type=users --author=test-author --template-name=exoservice-ls --model-name=user --description=testing --protection-level=public" in this application's directory
    Then it exits with code 1
    And I see the error "Service users already exists in this application"
