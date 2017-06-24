Feature: configuring the port

  As an Exosphere developer testing my application
  I want to be able to boot ExoCom at various ports
  So that I have flexibility in testing my system.

  - the default ExoCom port is 3100
  - provide the "EXOCOM_PORT" environment variable to listen for HTTP requests on a custom port


  Scenario: running at the default port
    When I run ExoCom
    Then it opens a port at 3100


  Scenario: the default port is already taken
    Given another service already uses port 3100
    When I try to run ExoCom
    Then it aborts with the message "port 3100 is already in use"


  Scenario: using a custom port
    When starting ExoCom at port 3200
    Then it opens a port at 3200


  Scenario: custom port that is already taken
    Given another service already uses port 3200
    When I try starting ExoCom at port 3200
    Then it aborts with the message "port 3200 is already in use"
