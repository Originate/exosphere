<table>
  <tr>
    <td><a href="02_web_server.md">&lt;&lt; the web server service</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="02_create_internal_service.md">creating an internal service &gt;&gt;</a></td>
  </tr>
</table>


# Scaffolding Services

<table>
  <tr>
    <td>
      Status: alpha - some things implemented but missing features
    </td>
  </tr>
</table>

Time to build the Todo service!
Let's first determine its API,
and then build it from the outside in.


## API

The data model for our Todo service is very simple:
* each person has one todo list
* a todo list contains many todo items
* a todo item is a small piece of text

The API for our Todo service is similarly simple:

<table>
  <tr>
    <th>command</th>
    <th>meaning</th>
    <th>outcome</th>
  </tr>
  <tr>
    <td>todo.list</td>
    <td>list all todo items for the given list</td>
    <td>todo.listing</td>
  </tr>
  <tr>
    <td>todo.add</td>
    <td>adds a todo item to the given list</td>
    <td>todo.addition</td>
  </tr>
</table>

A __command__ is a message send by one service to another
in order to make the latter service do something for the former.
It is usually written in present imperative.

An __outcome__ is the result of a command,
i.e. the message that is sent in response to it.
Usually those messages are a noun,
since they represent something that was created or is happening in response
to a command.


## Exoservices

To build our Todo service,
we could whip together another quick [ExpressJS](https://expressjs.com) stack
and use an [ExoRelay](https://github.com/Originate/exorelay-js) inside it,
as we did in the [web server service](02_web_server.md).
But for such a small micro-service a full web server stack
based on the MVC paradigm would be overkill:
* we don't need a sophisticated _routing_ layer,
  since micro-services don't serve a large variety of URLs:
  mostly, they just implement a single REST-like endpoint
* we don't need a sophisticated _model_ layer,
  since a service only deals with one model type
* we don't need _views_, since our service returns simple JSON data
  that can be automatically serialized
* since we don't have models and views,
  we don't really need _controllers_ either

All we need are simple handler functions
that are called for incoming messages.
They do a little bit of database work and send the outcome back.

The Exosphere SDK provides a number of tools
to build micro-services out of simple handler functions,
comparable to [AWS Lambda](https://aws.amazon.com/lambda).

Let's build our todo service using the
[Exoservice-JS](https://github.com/Originate/exoservice-js) framework.


## Building the Todo service

Let's have Exosphere generate the shell for our new service for us:

```
cd ~/todo-app
exo add exoservice-es6 todo-service list add
```

With this command we tell Exosphere:
Please add a micro-service written in EcmaScript 6,
named "todo-service", which responds to the commands "list" and "add"

Exosphere runs the
[Exoservice-ES6](https://github.com/Originate/exoservice-es6) generator,
which creates a fully functioning service shell for us.
Let's take a closer look.

The new service is in the directory `~/todo-app/todo-service`,
next to the web-server service we created in chapter 2-2.
It contains two files.

The service code is in __server.js__

```javascript
module.exports = {

  list: (_, {reply}) ->
    # TODO: do something in response to the "list" command here,
    #       and adjust the reply below
    reply 'list reply'

  add: (_, {reply}) ->
    # TODO: do something in response to the "add" command here
    #       and adjust the reply below
    reply 'add reply'

}
```

This file implements handlers for the three commands we told the generator about.
They don't do much except immediately sending a reply back.
We will add the business logic for them in a minute.

Next is the configuration file for our service.
It tells Exosphere about our service:

__config.yml__

```yaml
name: todo-service
description: Todo service

setup: npm install
startup:
  command: node_modules/.bin/exo-js
  online-text: online at port

messages:
  receives:
    - todo.add
    - todo.list
  sends:
    - todo.add-reply
    - todo.list-reply
```

Since our service is written in Node.JS,
our service directory also contains a very simple __package.json__ file:

```json
{
  "name": "todo-service",
  "description": "Todo service",
  "version": "0.0.1"
}
```

This file will later list any external NPM modules we use.

Exosphere encourages and supports a development model around best practices,
so the generator also creates a full suite of human-readable integration tests
in [Cucumber-JS](https://github.com/cucumber/cucumber-js).

Let's run them!

```
cd ~/todo-app/todo-service
spec
```

We see a bunch of passing tests that look like this:

```
....
```


## CLI for services

The `spec` command is part of a standardized set of CLI commands
that each service within Exosphere should define in its `bin` directory.
Since services in Exosphere can be written in any language and use any tool set,
a developer would have to learn a whole array of tools to install, boot, and test
the different services they use.
Thanks to this convention, you can use these commands across every service:

<table>
  <tr>
    <th>spec</th>
    <td>runs the full test suite</td>
  </tr>
  <tr>
    <th>spec <file name></th>
    <td>runs only the given test file</td>
  </tr>
  <tr>
    <th>lint</th>
    <td>runs all the linters</td>
  </tr>
  <tr>
    <th>build</th>
    <td>builds the service, i.e. compiles the source code into an executable</td>
  </tr>
  <tr>
    <th>watch</th>
    <td>
      builds continuously in the background,
      so that as you make changes to the source code,
      you always have a runnable version of the service
    </td>
  </tr>
  <tr>
    <th>setup</th>
    <td>installs all dependencies for the service, so that it is ready to be booted up</td>
  </tr>
  <tr>
    <th>run</th>
    <td>boots up the service</td>
  </tr>
</table>


## Building the business logic

Lets build the logic to save, list, and delete todo items into our service!
First, we extend the tests to reflect the behavior we would like our service to have.

> TODO: hardcore TDD action here

With a feature spec that describes the service properly in hand,
lets implement the service logic.

Since the data model that our Todo service has to manage is so simple,
we can use a simple (and very scalable) NoSQL technology to store our data.
Since todo items are some types of _documents_,
let's use MongoDB for this service.

> TODO: more hardcore TDD action

We now have a fully functioning todo service!

Lets integrate it into our application to see it in action!

> TODO: add "add todo" page to web server an check it out in action


Homework: add a "delete" function to the application,
so that you can keep your list of todo items relevant!


## Subscription Model

In the code above, we don't subscribe to messages,
but they still arrive at our service.

Services declare the messages they send and receive in their configuration file (config.yml).
When deployed, Exosphere reads this configuration and sets up the right
subscriptions and channels for each service.


Exosphere guarantees delivery of every message to at least one instance of
each subscribed service type.

Messages can be delivered to multiple instances,
for example if instances crash while processing the message, or fail to confirm
having received and processed it.
