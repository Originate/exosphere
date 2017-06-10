Feature: Manage new instances of services

  As an Exosphere operator
  I want that service instances can register themselves with ExoCom
  So that the system can add and remove instances at runtime and thereby scale with demand.

  Rules:
  - services send an "exocom.register" message when they come online
  - services are automatically de-registered when they go offline


  Scenario: a new service comes online
    Given an ExoCom instance configured with the routes:
      """
      [
        {
          "role": "users",
          "namespace": "foo"
        }
      ]
      """
    When a new "users" service instance registers itself with it via the message:
      | NAME    | exocom.register-service   |
      | PAYLOAD | { "clientName": "users" } |
    Then ExoCom now knows about these service instances:
      | CLIENT NAME | SERVICE TYPE | INTERNAL NAMESPACE |
      | users       | users        | foo                |


  Scenario: deregister a service once it goes offline
    Given an ExoCom instance configured with the routes
    """
    [
      {
      "role": "users",
      "namespace": "foo"
      },
      {
      "role": "tweets",
      "namespace": "bar"
      }
    ]
    """
    And a running "users" instance
    And a running "tweets" instance
    When the "tweets" service goes offline
    Then ExoCom now knows about these service instances:
      | CLIENT NAME | SERVICE TYPE | INTERNAL NAMESPACE |
      | users       | users        | foo                |
