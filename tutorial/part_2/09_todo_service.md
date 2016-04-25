<table>
  <tr>
    <td><a href="02_web_server.md">&lt;&lt; the web server service</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="02_create_internal_service.md">creating an internal service &gt;&gt;</a></td>
  </tr>
</table>


# Building the Todo service

<table>
  <tr>
    <td>
      <b><i>
        Status: beta - basics implemented, needs more hands-on testing
      </i></b>
    </td>
  </tr>
</table>

Time to build the Todo service!
Let's first determine its data model,
then based on that its API,
and then build the actual system from the outside in.


## Data Model

The data model for our Todo service is very simple:
* each person has one todo list
* a todo list contains many todo items
* a todo item is a small piece of text


## API

Such a simple data model can be accessed through a similarly simple API:

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
    <td>todo.added</td>
  </tr>
</table>


## Building the Todo service

Let's have Exosphere generate the shell for our new service for us:

```
cd ~/todo-app
exo add service
```

The command asks us a number of things it needs to know:

<table>
  <tr>
    <th>input field</th>
    <th>description</th>
    <th>we enter</th>
  </tr>
  <tr>
    <td>service name</td>
    <td>the name for the service to create</td>
    <td>todo-service</td>
  </tr>
  <tr>
    <td>language</td>
    <td>the language the service should be written in</td>
    <td>ecmascript</td>
  </tr>
  <tr>
    <td>incoming commands</td>
    <td>the incoming commands the service should listen to</td>
    <td>list, add</td>
  </tr>
  <tr>
    <td>outgoing messages</td>
    <td>the commands the service should send out</td>
    <td>listing, added</td>
  </tr>
</table>

Exosphere runs the
[Exoservice-ES6](https://github.com/Originate/exoservice-es6) generator,
which creates a fully functioning service shell for us.
The new service is in the directory `~/todo-app/todo-service`,
next to the web-server service we created in [chapter 2-5](05_web_server.md).
This is because we told Exosphere to _add_ the service to the current application
by running `exo add <service name>`.

The service code is in __server.js__

```javascript
module.exports = {

  "list": function listHandler (_, {reply}) {
    // TODO: handle the "list" command and adjust the reply below
    reply('list reply');
  },

  "add": function addHandler (_, {reply}) {
    // TODO: handle the "add" command and adjust the reply below
    reply('add reply');
  }

};
```

This file implements handlers for the two commands we told the generator about.
They don't do much except immediately sending a reply back.
We will add the business logic for them in a minute.

The other important file is the service configuration file.
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
    - todo.added
    - todo.listing
```

Since our service is written in Node.JS,
our service directory also contains a very simple __package.json__ file,
which is required by Node's package management system:

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
Let's run the tests!

```
cd ~/todo-app/todo-service
spec
```

We see a bunch of passing tests that look like this:

```
....
```

The generator has created a fully running service for us!
All we have to do is add its business logic and we are off to the races.


## Idempotency

Messages can get lost,
causing the sender to assume the transaction didn't happen
and re-sending a command.
Services must be able to recognize and deal with this situation.
For example, transactions that changes state
should keep a history of the recently performed commands
and check it for matches before executing an incoming command.


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


The todo-service now fully works.
Next we are going to integrate it into the web server.


<table>
  <tr>
    <td><a href="10_integration_into_web_server.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>

