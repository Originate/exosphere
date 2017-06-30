Feature: Translating messages

  As an ExoService developer
  I want to be able to use generic off-the-shelf services with a domain-specific API
  So that I can quickly plug together application prototypes without copy-and-pasting code around.

  - when a broadcasted message is delivered to a service,
    the service name in the message name is replaced with the service type
  - when a service broadcasts a message,
    Exocom replaces the service type with the external service name


  Background:
    Given an ExoCom instance configured with the routes:
    """
    [
      {
        "role": "web",
        "receives": ["tweets.created"],
        "sends": ["tweets.create"]
      },
      {
        "role": "tweets",
        "receives": ["text-snippets.create"],
        "sends": ["text-snippets.created"],
        "namespace": "text-snippets"
      }
    ]
    """
    And a running "web" instance
    And a running "tweets" instance


  Scenario: translating a message
    When the "web" service sends "tweets.create"
    Then ExoCom signals "web  --[ tweets.create ]-[ text-snippets.create ]->  tweets"
    And ExoCom broadcasts the message "text-snippets.create" to the "tweets" service


  Scenario: translating a reply
    When the "web" service sends "tweets.create"
    And the "tweets" service sends "text-snippets.created" in reply to "111"
    Then ExoCom signals "tweets  --[ text-snippets.created ]-[ tweets.created ]->  web  (XX ms)"
    And ExoCom broadcasts the reply "tweets.created" to the "web" service
