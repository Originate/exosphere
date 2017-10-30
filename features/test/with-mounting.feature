Feature: Mounting service directories in docker

  Scenario: can test mounted services
    Given I am in the root directory of the "frontend-with-webpack" example application
    And starting "exo test" in my application directory
    Then I eventually see the following snippets:
      | Testing service 'frontend-service'  |
      | All tests passed                    |
    And it exits with code 0
