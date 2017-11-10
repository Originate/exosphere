Feature: attempting to add a service without providing any templates

  As a developer adding features to an Exosphere application
  I want Exosphere to alert me when attempting to create a service without providing templates
  So that I know that ".exosphere/service_templates" folder either does not contain valid boilr templates or is empty


  Scenario: adding a service without providing any templates
    Given I am in the root directory of an empty application called "test-app"
    And my application contains the empty directory ".exosphere/service_templates"
    When starting "exo add" in my application directory
    Then I eventually see:
      """
      no templates found
      """
    And I eventually see:
      """
      Please add templates to the ".exosphere/service_templates" folder of your code base.
      """
    And it exits with code 1
