# Mock implementation of ExoCom in JavaScript

[![Dependency Status](https://david-dm.org/originate/exocom-mock-js.svg)](https://david-dm.org/originate/exocom-mock-js)
[![devDependency Status](https://david-dm.org/originate/exocom-mock-js/dev-status.svg)](https://david-dm.org/originate/exocom-mock-js#info=devDependencies)

> a mock implementation of [ExoCom-Dev](https://github.com/Originate/exocom-dev)
for sending and receiving messages to your ExoServices in test


## Installation

```
$ npm i --save-dev exocom-mock
```


## Usage

* create an instance

  ```coffeescript
  ExoComMock = require('exocom-mock')
  exocom = new ExoComMock
  ```

* register a service to send messages to

  ```coffeescript
  exocom.registerService name: 'users', port: 4001
  ```

* send a message to the service

  ```coffeescript
  exocom.sendMessage service: 'users', name: 'users.create', message-id: '123', payload: { name: 'Jean-Luc Picard' }
  ```

* verifying messages sent out by the service under test

  ```coffeescript
  # ... make your service sent out a request here via exocom.sendMessage...

  # wait for the message to arrive
  exocom.waitUntilReceive =>

    # verify the received message
    expect(exocom.receivedMessages()).to.eql [
      {
        name: 'users.created'
        payload:
          name: 'Jean-Luc Picard'
      }
    ]
  ```

* if you want to verify more received messages later, you can reset the register of received messages so far

  ```coffeescript
  exocom.reset()
  ```

* finally, close your instance when you are done, so that you can create a fresh one for your next test

  ```coffeescript
  exocom.close()
  ```


## Development

See our [developer documentation](CONTRIBUTING.md)
