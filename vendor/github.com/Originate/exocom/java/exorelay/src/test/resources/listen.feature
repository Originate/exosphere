Feature: Listening

  As a developer building Exosphere applications
  I want to be able to add an Exosphere communication relay to any code base
  So that I can write Exosphere services without constraints on my code layout.

  Rules
  - call "listen" on an ExoRelay instance to take it online
  - you provide the port as an argument to "listen"


  Scenario: Setting up at the given port
    Given ExoCom runs at port 4100
    And an ExoRelay instance
    When I take it online at port 4001
    Then it is online at port 4001
