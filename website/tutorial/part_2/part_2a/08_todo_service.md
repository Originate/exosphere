<table>
  <tr>
    <td><a href="07_exoservices.md">&lt;&lt; exoservice</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="09_integration_with_web_server.md">integration with the web server &gt;&gt;</a></td>
  </tr>
</table>


# Building the Todo service

Time to build the Todo service!
Storing todo entries is pretty straightforward.
We can start off with a basic service
that provides persistence via a
[CRUD](https://en.wikipedia.org/wiki/Create,_read,_update_and_delete)
API.
To keep things simple here,
we'll use [MongoDB](https://www.mongodb.com) for persistence.

Let's have the Exosphere CLI generate a fully functional service for us:

<a class="runMarkdown_consoleWithInputFromTable">
```
cd ~/todo-app
exo add service
```

The command asks us a number of things:

<table>
  <tr>
    <th>input field</th>
    <th>description</th>
    <th>we enter</th>
  </tr>
  <tr>
    <td>name</td>
    <td>the name for the service to create</td>
    <td>todo</td>
  </tr>
  <tr>
    <td>description</td>
    <td>description for the service to create</td>
    <td>stores the todo entries</td>
  </tr>
  <tr>
    <td>type</td>
    <td>the type of service we want to create</td>
    <td>exoservice-es6-mongodb</td>
  </tr>
</table>

</a>

The new service is in the directory `~/todo-app/todo`,
next to the web-server service we created in [chapter 2-4](04_html_server.md).


## Service configuration

Let's first check out the service configuration file.

<a class="runMarkdown_verifyFileContent">
__~/todo-app/todo/config.yml__

```yaml
name: todo
description: stores the todo entries

setup: npm install
startup:
  command: node node_modules/exoservice/bin/exo-js
  online-text: online at port
tests: node_modules/cucumber/bin/cucumber.js

messages:
  receives:
    - todo.create
    - todo.create_many
    - todo.delete
    - todo.list
    - todo.read
    - todo.update
  sends:
    - todo.created
    - todo.created_many
    - todo.deleted
    - todo.listing
    - todo.details
    - todo.updated
```
</a>

The __setup__ block describes how to make this service runnable.
The __startup__ block describes how to boot it up,
and determine whether the service is running and ready to accept traffic.
In this case it is told to wait for "online at port" as the console output.
The __messages__ block defines all the messages sent and received by this service.
This services provides a simple CRUD interface,
i.e. allows to created, read, update, and delete todo items.


## The service source code

The service code is in __~/todo-app/todo/src/server.js__.
It defines the handler functions for incoming messages.
Simplified (without handling errors and edge cases) it looks like this:

```javascript
module.exports = {

  'todo.create': (todoData, {reply}) => {
    collection.insertOne(todoData, (err, result) => {
      reply('todo.created', result)
    })
  },

  'todo.read': (query, {reply}) => {
    collection.find(query, (entry) => {
      reply('todo.details', entry)
    }))
  },

  'todo.update': (todoData, {reply}) => {
    collection.updateOne({ _id: todoData.id }, {$set: todoData}, (result) => {
      collection.find({_id: todoData.id}, (todo) => {
        reply('todo.updated', todo)
      })
    })
  },

  'todo.delete': (query, {reply}) => {
    collection.deleteOne({_id: id}, (result) => {
      reply('todo.deleted', query)
    })
  },

  'todo.list': (_, {reply}) => {
    collection.find({}, (todos) => {
      reply('todo.listing', todos)
    })
  }

};
```

This file is executed by [Exoservice-JS](https://github.com/originate/exoservice-js),
a tool that makes it extremely easy to build micro services
out of serverless lambda functions.


## Tests

Exosphere encourages best practices,
so this generator also creates a full suite of human-readable integration tests
using [Cucumber-JS](https://github.com/cucumber/cucumber-js)<sup>1</sup>.
Let's look at one:

__~/todo-app/todo/features/create.feature__

```cucumber
Feature: Creating todos

  Rules:
  - when successful, the service replies with "todo.created" and the created record
  - when there is an error, the service replies with "todo.not-created" and a message describing the error


  Background:
    Given an ExoCom server
    And an instance of this service


  Scenario: creating a valid todo record
    When receiving the message "todo.create" with the payload:
      """
      { "name": "Jean-Luc Picard" }
      """
    Then this service replies with "todo.created" and the payload:
      """
      {
        "id": /\d+/,
        "name": 'one'
      }
      """
    And the service now contains the todos:
      | NAME |
      | one  |
```

Since this is generated code,
the code examples aren't specific to todo entries,
but enough to get us started here.

Before our new service can do anything,
we need to get it ready
by installing its external dependencies:

<a class="runMarkdown_consoleWithDollarPrompt">
```
$ exo setup
```

</a>
With that out of the way, let's run the tests:

<a class="runMarkdown_consoleWithDollarPrompt">
```
$ exo test
```

</a>
The command prints a bunch of passing tests.
The generator has created a fully running service for us!


Takeaway:
> Exosphere makes it very easy to work on services.
> They can be scaffolded in a fully working form,
> and be worked on in isolation via test-driven development.

You can find a fully working version of the code base in your current state [here](code_09/todo-app).
Next we are going to integrate the todo service with the web server.


<table>
  <tr>
    <td><a href="09_integration_with_web_server.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>


<hr>

<sup>1</sup>
If you don't like [Cucumber](http://cucumber.io) or [ES6](http://es6-features.org/),
feel free to use another generator that uses other frameworks and languages,
or create your own!
