<table>
  <tr>
    <td><a href="07_exoservices.md">&lt;&lt; exoservice</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="09_integration_with_web_server.md">integration with the web server &gt;&gt;</a></td>
  </tr>
</table>


# Building the Todo service

> Note: This is the second service of the tutorial. 
If you have not created the first service, the `html-server`, 
you should [go back and create that service first](04_html_server.md). 
Otherwise, these instructions may not be accurate or 
the service may not work as expected.

Time to build the Todo service!
Storing todo entries is pretty straightforward.
We can start off with a basic service
that provides persistence via a
[CRUD](https://en.wikipedia.org/wiki/Create,_read,_update_and_delete)
API.
To keep things simple here,
we'll use [MongoDB](https://www.mongodb.com) for persistence.

Let's add the MongoDB service template and have the Exosphere CLI generate a fully 
functional service for us:

<a class="runMarkdown_consoleWithInputFromTable">

```
cd ~/todo-app
git clone https://github.com/originate/exoservice-js-mongodb .exosphere/service_templates/js-mongodb
exo add service
```

The command asks us a number of things:

<table>
  <tr>
    <th>prompt</th>
    <th>text you enter</th>
  </tr>
  <tr>
    <td>Choose a template</td>
    <td>2</td>
  </tr>
  <tr>
    <td>Service role</td>
    <td>todo</td>
  </tr>
  <tr>
    <td>Service type</td>
    <td>worker</td>
  </tr>
  <tr>
    <td>Description</td>
    <td>Stores the todo entries</td>
  </tr>
  <tr>
    <td>Model name</td>
    <td>todo</td>
  </tr>
  <tr>
    <td>author</td>
    <td>(Press enter and accept the default)</td>
  </tr>
  <tr>
    <td>data path</td>
    <td>(Press enter and accept the default)</td>
  </tr>
</table>

</a>

The new service is in the directory `~/todo-app/todo`,
next to the web-server service we created in [chapter 2-4](04_html_server.md).


## Service configuration

Let's first check out the service configuration file.

<a class="runMarkdown_verifyFileContent">
__~/todo-app/todo/service.yml__

```yaml
type: worker
description: stores the todo entries
author: 

startup:
  command: node src/server.js
  online-text: online at port

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

dependencies:
  - name: 'mongo'
    version: '3.4.0'
    config:
      volumes:
        - '{{EXO_DATA_PATH}}:/data/db'
      ports:
        - '27017:27017'
      online-text: 'waiting for connections'

development:
  scripts:
    run: node src/server.js
    test: node_modules/cucumber/bin/cucumber.js
```
</a>

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
Let's run the tests:

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
