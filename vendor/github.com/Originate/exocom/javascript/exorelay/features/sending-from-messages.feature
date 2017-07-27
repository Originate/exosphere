Feature: Sending from messages

  As an Exosphere developer
  I want to have access to a send method within my messages
  So that sending from messages is easy and my code concise.

  Rules:
  - call the "send" argument given to your message handler
    to send a message


  Background:
    Given ExoCom runs at port 4100
    And an ExoRelay instance called "exo-relay" running inside the "test-service" service


  Scenario: sending a message from within a message handler
    Given the "users.login" message has this handler:
      """
      exo-relay.register-handler 'users.create', (_payload, {send}) ->
        send 'passwords.encrypt', 'secret'
      """
    When receiving this message:
      """
      name: 'users.create'
      id: '123'
      """
    Then my message handler sends out the message:
      """
      name: 'passwords.encrypt'
      sender: 'test-service'
      payload: 'secret'
      id: '<%= request_uuid %>'
      """

