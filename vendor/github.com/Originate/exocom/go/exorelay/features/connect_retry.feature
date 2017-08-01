Feature: Retrying

  When running Exocom applications
  I want ExoRelays to reconnect after network failures
  So that my applications is resilient to Exocom outages.

  Rules
  - ExoRelay will retry connecting to Exocom with exponential backoff

  Background:
    Given an ExoRelay with the role "foo"

  Scenario: ExoRelay attempts to establish a connection to ExoCom
    Given Exocom is offline
    When ExoRelay and Exocom boot up simultaneously
    Then ExoRelay should connect to Exocom

  Scenario: ExoRelay will attempt to reconnect to ExoCom if ExoCom goes offline
    Given an ExoRelay instance that is connected to Exocom
    When Exocom goes offline momentarily
    Then ExoRelay should reconnect to Exocom
