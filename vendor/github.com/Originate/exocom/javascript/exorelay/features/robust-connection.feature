Feature: resilient connecting to ExoCom

  As an Exosphere operator
  I want ExoRelay to keep trying to connect to ExoCom
  So that I can deploy ExoCom and Exoservices in any order.

  Rules
  - failed connections to ExoCom are retried every second until a connection is made

  Scenario: connect to ExoCom instance that comes online after ExoRelay
    Given a new ExoRelay instance connecting to port 4100
    When an ExoCom instance comes online at port 4100, 1 second later
    Then ExoRelay connects to ExoCom

  Scenario: reconnect to ExoCom instance that goes down
    Given ExoCom runs at port 4000
    And an ExoRelay instance
    When ExoCom goes down for 1 second and comes back online
    Then ExoRelay is connected to ExoCom
