# Communication

Exosphere uses fully asynchronous communication that is decoupled in
* space: the sender does not need to know the address of the receiver.
  It simply sends out the message and trusts that the right receivers will listen to it
* time: the receiver doesn't have to be online or ready to receive when the message is sent.
  The message waits in a queue until it is processed.

This means messages are fire-and-forget.
You send off a message, and if there is a response, you will receive it
as a separate message.


## Message Protocol

To talk to a service, send off a message.

Messages have the following structure:
* __name:__ a random string like `users.create`
* __version:__ version of the message
* __payload:__ JSON data payload
* __id:__ id of the message
* __response-to:__ id of the messages that this message is a response to
* __sender:__ name of the service that sent the message

The structure of the payload can be defined using [JSON schema](http://json-schema.org).

Exosphere guarantees delivery of every message to at least one instance of
each subscribed service type.

Messages can be delivered to multiple instances,
for example if instances crash while processing the message, or fail to confirm
having received and processed it.


## Subscription Model

Services declare the messages they send and receive in their configuration file (config.yml).
When deployed, Exosphere reads this configuration and sets up the right
subscriptions and channels for each service.
