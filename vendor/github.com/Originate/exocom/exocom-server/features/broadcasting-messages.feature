Feature: Broadcasting messages

  As an ExoService
  I want to broadcast messages to other services
  So that I can communicate with them.

  - to broadcast a message, send a WebSocket request to the WebSocket port of the ExoCom Instance
  - ExoCom sends the message to all subscribed services


  Background:
    Given an ExoCom instance configured with the routes:
    """
    [
      {
        "role": "web",
        "receives": ["users.created"],
        "sends": ["users.create", "users.list", "delete user"]
      },
      {
        "role": "users",
        "receives": ["mongo.create", "delete user"],
        "sends": ["mongo.created"],
        "namespace": "mongo"
      },
      {
        "role": "log",
        "receives": ["users.create"]
      }
    ]
    """
    And a running "web" instance
    And a running "users" instance
    And a running "log" instance


  Scenario: broadcasting a message
    When the "web" service sends "users.create"
    Then ExoCom signals "web  --[ users.create ]-[ mongo.create ]->  users"
    And ExoCom broadcasts the message "mongo.create" to the "users" service


  Scenario: broadcasting a reply
    When the "web" service sends "users.create"
    And the "users" service sends "mongo.created" in reply to "111"
    Then ExoCom signals "users  --[ mongo.created ]-[ users.created ]->  web  (XX ms)"
    And ExoCom broadcasts the reply "users.created" to the "web" service


  Scenario: broadcasting a message with whitespace
    When the "web" service sends "delete user"
    Then ExoCom signals "web  --[ delete user ]->  users"


  # ERROR HANDLING
  Scenario: broadcasting an invalid message
    When the "web" service sends "users.get-SSN"
    Then ExoCom signals the error "Service 'web' is not allowed to broadcast the message 'users.get-SSN'"


  Scenario: broadcasting an invalid message
    When the "log" service sends "users.get-SSN"
    Then ExoCom signals the error "Service 'log' is not allowed to broadcast the message 'users.get-SSN'"


  Scenario: broadcasting a message with no receivers
    When the "web" service sends "users.list"
    Then ExoCom signals the error "Warning: No receivers for message 'users.list' registered"
