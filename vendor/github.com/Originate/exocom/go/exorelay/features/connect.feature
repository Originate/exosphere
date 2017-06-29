Feature: Connecting

  As a developer building Exosphere applications
  I want to be able to add an Exosphere communication relay to any code base
  So that I can write Exosphere services without constraints on my code layout.

  Rules
  - call "Connect" on an ExoRelay instance to take it online


  Background:
    Given an ExoRelay with the role "foo"


  Scenario: Setting up the ExoRelay instance
    When ExoRelay connects to Exocom
    Then it registers by sending the message "exocom.register-service" with payload:
      """
      { "clientName": "foo" }
      """
