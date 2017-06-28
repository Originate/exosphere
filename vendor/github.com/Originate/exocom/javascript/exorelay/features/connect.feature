Feature: Connecting

  As a developer building Exosphere applications
  I want to be able to add an Exosphere communication relay to any code base
  So that I can write Exosphere services without constraints on my code layout.

  Rules
  - call "connect" on an ExoRelay instance to take it online


  Background:
    Given ExoCom runs at port 4100


  Scenario: Setting up the ExoRelay instance
    Given a new ExoRelay instance connecting to port 4100
    When I take it online
    Then it connects to the given ExoCom host and port


  Scenario: Attempting to set up ExoRelay at an invalid port
    Given a new ExoRelay instance connecting to port 4000
    When I try to take it online
    Then ExoRelay emits an "error" event with the error "connect ECONNREFUSED 127.0.0.1:4000"
