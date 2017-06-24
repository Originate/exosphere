Feature: Defining the port at which the server listens

  As an ExoService developer
  I want to be able to run my service at a configurable port
  So that my services fit seamlessly into the networking setup of my infrastructure.


  Rules:
  - call "exo-js run" in the directory of an Exosphere service to start up that service
  - provide the ExoCom address via the environment variables EXOCOM_HOST and EXOCOM_PORT
  - provide the service name via the environment variable ROLE


  Background:
    Given an ExoCom instance running at port 3001

  Scenario: an ExoService instance connects to ExoCom
    When starting a service configured for ExoCom port 3001
    Then it connects to the ExoCom instance
