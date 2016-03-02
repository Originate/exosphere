# Creating Micro-Services

ExoSphere applications break up large (and unmaintainable) complexity
into many small, maintainable, and reusable services.
These services are completely independent code bases.
They run stand-alone, and can (hopefully) be reused in different applications.


## Servicifying an existing code base

Any code base can act as an Exoservice,
as long as it is able to send and receive Exosphere messages.
Implementing this functionality is easy,
since messages are transmitted via [simple HTTP POST requests](wire-format.md).

To make this even easier,
the Exosphere SDK provides libraries called __communication relays__
that provide those communication facilities for all mainstream stacks:
* [ExoRelay-JS](https://github.com/Originate/exorelay-js) for Node.js
  code bases.


## Lambda Services

Lamba services are the easiest way to implement microservices.
They contain of 2 files:
* a configuration file that tells Exosphere about this service
* a source code file that provides handlers for all incoming message types

Exosphere provides frameworks to write lambda services in most popular languages:
* [exoservice-js](https://github.com/Originate/exoservice-js) for services in Node.js

Example lambda services:
* [users service](https://github.com/Originate/exosphere-users-service)
