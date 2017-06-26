Feature: status check

  As an administrator
  I want to have an easy way to check whether an ExoRelay is online
  So that I know whether my system is set up correctly.

  Rules:
  - send a 'ping' message to ExoRelay to check that it is online
  - if it is, it returns a 'pong' response
  - any other return code or non-response means an error


  Scenario: ExoRelay is online
    Given ExoCom runs at port 4100
    And an ExoRelay instance listening on port 4000
    When I check the status
    Then it signals it is online
