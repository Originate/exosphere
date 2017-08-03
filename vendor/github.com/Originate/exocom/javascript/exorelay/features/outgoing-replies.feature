Feature: Sending outgoing replies to incoming messages

  As an Exosphere developer
  I want my code to be able to send replies to incoming Exosphere messages
  So that I can write responsive ExoServices.

  Rules:
  - call the "reply" argument given to your message handler
    to send a reply to the message you are currently processing


  Background:
    Given ExoCom runs at port 4100
    And an ExoRelay instance called "exo-relay" running inside the "test-service" service


  Scenario: sending a reply with JSON data
    Given the "users.create" message has this handler:
      """
      exo-relay.register-handler 'users.create', (user-attributes, {reply}) ->
        # on this line we would create a user record with the given attributes in the database
        reply 'users.created', id: 456, name: user-attributes.name
      """
    When receiving this message:
      """
      name: 'users.create'
      payload:
        name: 'Will Riker'
      id: '123'
      """
    Then my message handler replies with the message:
      """
      name: 'users.created'
      sender: 'test-service'
      payload:
        id: 456
        name: 'Will Riker'
      id: '<%= request_uuid %>'
      response-to: '123'
      """


  Scenario: sending a reply with string payload
    Given the "users.create" message has this handler:
      """
      exo-relay.register-handler 'ping', (_payload, {reply}) ->
        reply 'pong', 'from the test-service'
      """
    When receiving this message:
      """
      name: 'ping'
      id: '123'
      """
    Then my message handler replies with the message:
      """
      name: 'pong'
      sender: 'test-service'
      payload: 'from the test-service'
      id: '<%= request_uuid %>'
      response-to: '123'
      """
