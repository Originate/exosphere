Feature: test templates

  When developing an exosphere service template
	I want to be able to test it to make sure it can be used by exocom
	So that I don't have broken templates

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

  Scenario: success
    Given I am in the root directory of the "good" example template
    When running "exo template test" in my template directory
    Then it prints "Template passes" in the terminal
    And it exits with code 0
