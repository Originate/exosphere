Feature: Registering Exoservices


  Scenario: a service registers
    Given an ExoComMock instance
    When a new service instance registers itself with it via the message:
      | NAME    | exocom.register-service |
      | PAYLOAD | { "name": "users" }     |
    Then ExoCom now knows about these services:
      | users |
