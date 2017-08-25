Feature: test templates

  As a developer writing exsophere service templates
  I want to be able to see helpful messages when my templates are not valid
  So that I don't have broken templates

  Rules:
  - run "exo template test" in the directory of an exosphere template to test
    adding it to an application and run the tests

  - Structure rules for a valid template
    - contains "project.json"
    - contains "template" folder
      - contains a single folder
        - contains "service.yml"
        - contains "Dockerfile"
        - contains "tests/Dockerfile"

  - Template has default values for all required fields so running "exo add"
    with the template and entering nothing for all prompts does not fail

  - Running "exo test" in the generated service passes

  Scenario: missing project.json
    Given I am in the root directory of the "missing_project_json" example template
    When starting "exo template test" in my template directory
    Then I see:
      """
      Missing file: 'project.json'
      """
    And it exits with code 1

  Scenario: missing template folder
    Given I am in the root directory of the "missing_template_folder" example template
    When starting "exo template test" in my template directory
    Then I see:
      """
      Missing directory: 'template'
      """
    And it exits with code 1

  Scenario: missing service folder
    Given I am in the root directory of the "missing_service_folder" example template
    When starting "exo template test" in my template directory
    Then I see:
      """
      'template' directory must contain a single subdirectory
      """
    And it exits with code 1

  Scenario: missing Dockerfile
    Given I am in the root directory of the "missing_dockerfile" example template
    When starting "exo template test" in my template directory
    Then I see:
      """
      Missing file: 'template/{{serviceRole}}/Dockerfile'
      """
    And it exits with code 1

  Scenario: missing service.yml
    Given I am in the root directory of the "missing_service_yml" example template
    When starting "exo template test" in my template directory
    Then I see:
      """
      Missing file: 'template/{{serviceRole}}/service.yml'
      """
    And it exits with code 1

  Scenario: missing test Dockerfile
    Given I am in the root directory of the "missing_test_dockerfile" example template
    When starting "exo template test" in my template directory
    Then I see:
      """
      Missing file: 'template/{{serviceRole}}/tests/Dockerfile'
      """
    And it exits with code 1

  Scenario: test failure
    Given I am in the root directory of the "test_failure" example template
    When starting "exo template test" in my template directory
    Then it prints "Tests failed" in the terminal
    And it prints "running tests ... failing!" in the terminal
    And it exits with code 1
