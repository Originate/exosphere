Feature: Service endpoint enviornment variables

  As an Exosphere developer
  I want the services to know where each other are available
  So that communication between services are possible

  Rules:
  - "<NAME>_EXTERNAL_ORIGIN" is passed as an environment variable to
    each service for any public service
  - "<NAME>_INTERNAL_ORIGIN" is passed as an environment variable to
    each service for any public service


  Scenario: external origin env var available at run time
    Given I am in the root directory of the "service-external-origin-env" example application
    And starting "exo run" in my application directory
    And it prints "Listening on port 3000" in the terminal
    And it prints "webpack: Compiled successfully" in the terminal
    Then http://localhost:3010 displays:
      """
      Backend located at http://localhost:3000
      """

  Scenario: internal origin env var available at run time (public service)
    Given I am in the root directory of the "service-internal-origin-public" example application
    And starting "exo run" in my application directory
    And it prints "frontend service online" in the terminal
    And it prints "backend service online" in the terminal
    Then http://localhost:3010 displays:
      """
      Backend service content
      """

  Scenario: internal origin env var available at run time (private service)
    Given I am in the root directory of the "service-internal-origin-private" example application
    And starting "exo run" in my application directory
    And it prints "frontend service online" in the terminal
    And it prints "backend service online" in the terminal
    Then http://localhost:3000 displays:
      """
      Backend service content
      """
