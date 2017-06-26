# Exocom Development Guidelines


To boot up Exocom by itself on the command line.

```
$ exocom [--port <port number>]
```

The default port is 3100.
More details around configuring the port [here](features/configuring-the-ports.feature).


## Architecture

ExoCom is implemented using the _micro-kernel_ pattern:
as a number of relatively independent subsystems that are integrated through
a small and lightweight core that only provides
communication between the subsystems.
Each subsystem provides a particular set of functionality,
and is tested and implemented independently,
as a set of one or more classes.

* [HttpListener](src/http-listener):
  implements the HTTP endpoint that services talk to in order to make requests to ExoCom
* [ClientRegistry](src/client-registry):
  keeps track of which service is allowed to send and receive what messages
* [MessageSender](src/message-sender):
  sends messages to external services
* [Kernel](src/exocom.ls):
  integrates all the above subsystems and provides the
  [programmatic API](#javascript-api)

<img src="documentation/architecture.gif">


## Configuring the routes

To provide the routing information to Exocom,
do a POST request to `http://localhost:<exocom-port>/services`
with the payload described [here](features/configuring-services.feature).


## Testing

The tests run against the compiled output, so you need to run `watch` before executing them.

- run all tests: `spec`
- run feature tests: `features`
  - set environment variable `DEBUG_EXOCOM_SERVER=1` to print debug output
- run unit tests: `tests`
- run linter: `lint`


## Publishing

* update dependencies: `update`
* publish a new version: `publish <patch|minor|major>`
