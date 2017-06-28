Feature: Configuring the ExoCom address

  As a developer
  I want to be able to configure the ExoCom address that my ExoRelay instance is talking to
  So that I have flexibility in my test setup.

  Rules:
  - provide the ExoCom host via the "exocomHost" constructor parameter
  - provide the ExoCom port via the "exocomPort" constructor parameter


  Scenario: the user provides the ExoCom address
    When creating an ExoRelay instance using ExoCom host "localhost" and port 4100
    Then this instance uses the ExoCom host "localhost" and port 4100


  Scenario: the user does not provide the ExoCom port
    When trying to create an ExoRelay without providing the ExoCom port
    Then it throws the error "ExoCom port not provided to Exorelay"


  Scenario: the user does not provide the ExoCom host
    When trying to create an ExoRelay without providing the ExoCom host
    Then it throws the error "ExoCom host not provided to Exorelay"
