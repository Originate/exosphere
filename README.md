# Exosphere SDK
[![Circle CI](https://circleci.com/gh/Originate/exosphere-sdk.svg?style=shield&circle-token=fc8148ed828cc81e6ca44920672af8f773106795)](https://circleci.com/gh/Originate/exosphere-sdk)

This is the command-line interface of the Exosphere SDK.


## Architecture

This SDK consists of the following components:

* [exo-install command](commands/install):
  Sets up a freshly cloned Exosphere application
  by running the setup scripts for each service.

* [exo-run command](commands/run):
  Starts a properly set up (using `exo-install`) Exosphere application
  by launching all of its services.

* [feature specs](features): written in Cucumber

* [ExoComm-dev](https://github.com/Originate/exocomm-dev):
  development version of the Exosphere communication infrastructure


## How to use

* install this SDK on your machine

  ```
  npm i -g exosphere-sdk
  ```

* download an Exosphere application, for example
  [SpaceTweet](https://github.com/Originate/exosphere--example-app--space-tweet)

* set up the application:

  ```
  cd <application folder>
  exo-install
  ```

* start the application:

  ```
  exo-run
  ```


## Message protocol

The message protocol used in this prototype
is optimized for extreme simplicity.
This is to encourage experimentation
and implementations using many different languages.

* ExoComm sends messages to your service
  via an HTTP POST request to

  ```url
  /run/<message name>
  ```

* the request body
  is a JSON data structure
  with the format:

  ```
  {
    requestId: <the ID of this message, as a UUID>
    payload: <optional, a string, array, hash, or null>
    responseTo: <optional, the requestId of the command responded to
  }
  ```

* your service
  can send replies or new outgoing messages
  via an HTTP POST request to

  ```url
  http://localhost:<exocomm-port>/send/<message-name>
  ```

* the request body is like in incoming requests


## Build your own service

Converting an existing service into an Exoservice requires only two steps:
* listen at the port provided via the environment variable `EXORELAY_PORT`
  for incoming requests to

  ```
  /run/<message name>
  ```

  and understand the messages sent to your service this way

* send out messages and replies via POST requests to `EXORELAY_URL`

* if you do this with a lot of services,
  try to use an (or build your own) service framework like
  [Exoservice-JS](https://github.com/Originate/exoservice-js)


## Development

see our [developer guidelines](CONTRIBUTING.md)
