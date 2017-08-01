Feature: configuring the port

  As an Exosphere developer testing my application
  I want to be able to boot ExoCom at various ports
  So that I have flexibility in testing my system.

  - the default ExoCom port is 3100
  - provide the "EXOCOM_PORT" environment variable to listen for HTTP requests on a custom port


  Scenario: running at the default port
    When starting ExoCom
    Then I see "ExoCom online at port 3100"


  Scenario: using a custom port
    When starting ExoCom at port 3200
    Then I see "ExoCom online at port 3200"


  Scenario:  port that is already taken
    Given another service already uses port 3300
    When trying to start ExoCom at port 3300
    Then it aborts with the message "address already in use"
