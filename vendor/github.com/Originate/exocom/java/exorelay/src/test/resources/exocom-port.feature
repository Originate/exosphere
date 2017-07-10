Feature: Configuring the ExoCom port

  As a developer
  I want to be able to configure the ExoCom port that my ExoRelay instance is talking to
  So that I have flexibility in my test setup.

  Rules:
  - the default ExoCom port is 4100
  - provide a custom ExoCom port via the "exocomPort" constructor parameter


  Scenario: the user does not provide the ExoCom port
    When I try to create an ExoRelay without providing the ExoCom port
    Then it throws the error "exocomPort not provided"


  Scenario: the user provides an available ExoCom port
    When I create an ExoRelay instance that uses ExoCom port 3200
    Then this instance uses the ExoCom port 3200
