# Exosphere Communication Relay for JavaScript

> Communication relay between JavaScript code bases and the Exosphere environment

[![Dependency Status](https://david-dm.org/originate/exorelay-js.svg)](https://david-dm.org/originate/exorelay-js)
[![devDependency Status](https://david-dm.org/originate/exorelay-js/dev-status.svg)](https://david-dm.org/originate/exorelay-js#info=devDependencies)

This library allows you to add Exosphere communication to any Node.JS codebase.
It is intended to be used in your web or API server,
or with legacy Node code bases.
If you want to write a new micro-service in Node,
please use [ExoService-JS](https://github.com/Originate/exoservice-js),
which uses this library internally.


## Add an ExoRelay to your application

Each code base should have only one ExoRelay instance.

```coffeescript
ExoRelay = require 'exorelay'

exoRelay = new ExoRelay exocomPort: <port>, serviceName: <name of the service using ExoRelay>
exoRelay.on 'online', (port) ->  # yay, we are online!
exoRelay.on 'error', (err) ->    # examine, print, or log the error here
exoRelay.listen 4000
```

More details and how to customize the port is described in the [spec](features/listen.feature).

## Events

ExoRelay instances are [EventEmitters](https://nodejs.org/api/events.html).
They emit the following events to signal state changes:

<table>
  <tr>
    <th>online</th>
    <td>The instance is completely online now. Provides the port it listens on.
  </tr>
  <tr>
    <th>offline</th>
    <td>The instance is offline now.</td>
  </tr>
  <tr>
    <th>error</th>
    <td>An error has occurred. The instance is in an invalid state, your application should crash.</td>
  </tr>
</table>


## Handle incoming messages

Let's say we build a service that greets users.
Here is how to register a handler for incoming "hello" messages:

```coffeescript
exoRelay.registerHandler 'hello', (name) -> console.log "hello #{name}!"
```

More details on how to define message listeners are [here](features/receiving-messages.feature).


## Send outgoing messages

Send a message to Exosphere:

```coffeescript
exoRelay.send 'hello', name: 'world'
```

Sending a message is fire-and-forget, i.e. you don't have to wait for the
sending process to finish before you can do the next thing.
More details on how to send various data are [here](features/sending.feature).


## Send outgoing replies to incoming messages

If you are implementing services, you want to send outgoing replies to incoming messages:

```coffeescript
exoRelay.registerHandler 'user.create', (userData, {reply}) ->
  # on this line we would save userData in the database
  reply 'user.created', id: 456, name: userData.name
```

More details and a working example of how to send replies is [here](features/outgoing-replies.feature).


## Handle incoming replies

If a message you send expects a reply,
you can provide the handler for it right when you send it:

```coffeescript
exoRelay.send 'users.create', name: 'Will Riker', (createdUser) ->
  console.log "the remove service has finished creating user #{createdUser.id}"
```

More examples for handling incoming replies are [here](features/incoming-replies.feature).


## Development

See our [developer guidelines](CONTRIBUTING.md)
