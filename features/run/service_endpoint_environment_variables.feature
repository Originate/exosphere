Feature: Service endpoint enviornment variable

  As an Exosphere developer
  I want the services to know where each other are available from the outside world
  So clients can talk directly to other services

  Rules:
  - "SERVICE_<NAME>_EXTERNAL_ORIGIN" is passed as an environment variable to
    each service for any public service


  Scenario: external origin env var available at run time
    Given I am in the root directory of the "service-exteral-origin-env" example application
    And starting "exo run" in my application directory
    And it prints "Listening on port 3000" in the terminal
    And it prints "webpack: Compiled successfully" in the terminal
    Then http://localhost:3010 displays:
      """
      Backend located at http://localhost:3000
      """

  Scenario: external origin env var available at build time
    Given I am in the root directory of the "service-exteral-origin-env" example application
    And starting "exo run --production" in my application directory
    And it prints "Listening on port 3000" in the terminal
    Then http://localhost:3010 displays:
      """
      Backend located at http://localhost:3000
      """
