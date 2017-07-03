Feature: Attempting to add a duplicate service

  As a developer adding features to an Exosphere application
  I want Exosphere to alert me when attempting to create a service without providing templates
  So that I know that .exosphere folder either does not contain valid boilr templates or is empty


  Scenario: Adding a service without providing any templates
    Given I am in the root directory of an empty application called "test app"
    And my application contains the empty directory ".exosphere"
    When starting "exo-add" in this application's directory
    Then it exits with code 1
    And I see the error "no templates found"
