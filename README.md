# Exosphere SDK
[![Circle CI](https://circleci.com/gh/Originate/exosphere-cli-run.svg?style=shield&circle-token=fc8148ed828cc81e6ca44920672af8f773106795)](https://circleci.com/gh/Originate/exosphere-cli-run)

This is the command-line interface of the Exosphere SDK.


## Architecture

This SDK consists of the following components:

* [ExoComm-dev](https://github.com/Originate/exocomm-dev):
  development version of the Exosphere communication infrastructure

* [exo-install command](commands/install):
  Sets up a freshly cloned Exosphere application
  by running the setup scripts for each service.

* [exo-run command](commands/run):
  Starts a properly set up (using `exo-install`) Exosphere application
  by launching all of its services.

* [feature specs](features): written in Cucumber


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


## Development

see our [developer guidelines](CONTRIBUTING.md)
