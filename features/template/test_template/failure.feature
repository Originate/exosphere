Feature: test templates

  When developing an exosphere service template
	I want to be able to see helpful messages when my templates are not valid
	So that I don't have broken templates

  - run "exo template test" in the directory of an exosphere template to test
    adding it to an application and run the tests

  Scenario: missing Dockerfile
    Given I am in the root directory of the "missing_dockerfile" example template
    When running "exo template test" in my template directory
    Then it prints "Template fails" in the terminal
    And it exits with code 1

  Scenario: missing service folder
    Given I am in the root directory of the "missing_service_folder" example template
    When running "exo template test" in my template directory
    Then it prints "Template fails" in the terminal
    And it exits with code 1

  Scenario: missing service.yml
    Given I am in the root directory of the "missing_service_yml" example template
    When running "exo template test" in my template directory
    Then it prints "Template fails" in the terminal
    And it exits with code 1

  Scenario: missing template folder
    Given I am in the root directory of the "missing_template_folder" example template
    When running "exo template test" in my template directory
    Then it prints "Template fails" in the terminal
    And it exits with code 1

  Scenario: missing test Dockerfile
    Given I am in the root directory of the "missing_test_dockerfile" example template
    When running "exo template test" in my template directory
    Then it prints "Template fails" in the terminal
    And it exits with code 1

  Scenario: test failure
    Given I am in the root directory of the "test_failure" example template
    When running "exo template test" in my template directory
    Then it prints "Template fails" in the terminal
    And it exits with code 1
