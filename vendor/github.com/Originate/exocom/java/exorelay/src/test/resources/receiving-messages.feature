Feature: Receiving messages

  As an Exosphere developer
  I want my code base to be able to respond to incoming messages
  So that I can create Exoservices.

  Rules:
  - you need to register handlers for messages that you want to receive
  - handlers are called with the request data

#
#  Background:



  Scenario: receiving a message without payload
    Given ExoCom runs at port 4100
    And an ExoRelay instance listening on port 4000
    Given I register a handler for the "hello-world" message:
      """
      exo-relay.register-handler('hello-world', param -> { System.out.println("Hello world!");});
      """
    When receiving this message:
      """
      {
      "name": "hello-world",
      "id": "123"
      }
      """
#    Then ExoRelay runs the registered handler, in this example calling print with "Hello world!"

#
#  Scenario: Receiving a message with string payload
#    Given I register a handler for the "hello" message:
#      """
#      exo-relay.register-handler 'hello', (name) ->
#        print "Hello #{name}!"
#      """
#    When receiving this message:
#      """
#      name: 'hello'
#      payload: 'world'
#      id: '123'
#      """
#    Then ExoRelay runs the registered handler, in this example calling "print" with "Hello world!"
#
#
#  Scenario: receiving a message with Hash payload
#    Given I register a handler for the "hello" message:
#      """
#      exo-relay.register-handler 'hello', ({name}) ->
#        print "Hello #{name}!"
#      """
#    When receiving this message:
#      """
#      name: 'hello'
#      payload:
#        name: 'world'
#      id: '123'
#      """
#    Then ExoRelay runs the registered handler, in this example calling "print" with "Hello world!"
#
#
#  Scenario: Receiving a message with array payload
#    Given I register a handler for the "sum" message:
#      """
#      exo-relay.register-handler 'sum', (numbers) ->
#        print numbers[0] + numbers[1]
#      """
#    When receiving this message:
#      """
#      name: 'sum'
#      payload: [1, 2]
#      id: '123'
#      """
#    Then ExoRelay runs the registered handler, in this example calling "print" with "3"
#
#
# #   ERROR CHECKING
#
#  Scenario: missing incoming message
#    When receiving this message:
#      """
#      name: undefined
#      id: '123'
#     """
#    Then ExoRelay emits an "error" event with the error "unknown message: 'undefined'"
#
#
#  Scenario: the incoming message is not registered
#    When receiving this message:
#      """
#      name: 'zonk'
#      id: '123'
#      """
#    Then ExoRelay emits an "error" event with the error "unknown message: 'zonk'"
