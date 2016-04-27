<table>
  <tr>
    <td><a href="09_exoservice.md">&lt;&lt; exoservice</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="11_integration_into_web_server.md">integration into the web server &gt;&gt;</a></td>
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
    <td>ES6</td>
  </tr>
  <tr>
    <td>datastore</td>
    <td>the backing data store, if any</td>
    <td>MongoDB</td>
  </tr>
  <tr>
    <td>incoming commands</td>
    <td>the incoming commands the service should listen to</td>
    <td>list</td>
  </tr>
  <tr>
    <td>outgoing messages</td>
    <td>the commands the service should send out</td>
    <td>listing</td>
  </tr>
</table>

Since the data model for our Todo service is so simple,
we can use a simple (and scalable) NoSQL technology to store our data.
Since todo items are some types of _documents_,
potentially with more attributes down the line ("checked"),
let's use MongoDB for storing them.

Exosphere runs the
[Exoservice-ES6](https://github.com/Originate/exoservice-es6) generator,
which creates a fully functioning service shell for us.
The new service is in the directory `~/todo-app/todo-service`,
next to the web-server service we created in [chapter 2-5](05_web_server.md).
This is because we told Exosphere to _add_ the service to the current application
by running `exo add <service name>`.

The service code is in __~/todo-app/todo-service/server.js__

```javascript
module.exports = {

  "list": function listHandler (_, {reply}) {
    // TODO: handle the "list" command and adjust the reply below
    reply('list reply');
  }

};
```

This file implements handlers for the two commands we told the generator about.
They don't do much except immediately sending a reply back.
We will add the business logic for them in a minute.

The other important file is the service configuration file.
It tells Exosphere everything it needs to know about our service:

__~/todo-app/todo-service/config.yml__

```yaml
name: todo-service
description: Todo service

startup:
  online:
    console-output: online at port

messages:
  receives:
    - todo.list
  sends:
    - todo.listing
```

The _startup_ block describes how Exosphere determines
whether the service has successfully booted up.
In this case it is told to wait for the console output "online at port".

Since our service is written in Node.JS,
our service directory also contains a very simple __package.json__ file,
which is required by Node's package management system:

__~/todo-app/todo-service/package.json__

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

Let's look at one:

__~/todo-app/todo-service/features/add__

```cucumber
Feature: "list" command

  Describe the intent of the add command here.


  Scenario: sending a "list" command
    Given a running "todo" service
    When receiving a "list" command with the payload:
      """
      {}
      """
    Then it replies with "listing" and the payload:
      """
      {}
      """
```

This describes exactly what our simple service shell does at this point!
Let's run the tests:

```
$ cd ~/todo-app/todo-service
$ spec
```

We see a bunch of passing tests.
The generator has created a fully running service for us!
All we have to do is fill in the business logic and we are off to the races!


## Building the business logic

Lets build the logic for our first command - listing todo items - into our service!
First we extend the tests to reflect the behavior we would like our service to have.

In good TDD fashion, we first think about what our service should do.
Here is a first attempt at a description of the `list` API of our todo service:

__~/todo-app/todo-service/features/list.feature__

```cucumber
Feature: Listing todo items

  As the web server
  I want to get all todo items of my current user
  So that I can display them on the homepage.


  Scenario: the user has no todo items
    Given a running "todo" service
    When receiving a "list" command with the payload:
      """
      {
        userId: 1
      }
      """
    Then it replies with "listing" and the payload:
      """
      []
      """
```

We have fleshed out the description more.
In particular, the specification for the "list" feature
now has a proper user story,
and our scenario has a proper title.
Let's run this test to see it fail:

```
$ spec
```

The issue is that our service returns the wrong data structure.
It returns a hash, but is now specified to return a list (array).
Let's change the return value on line XXX to `[]` and run our tests again.
Now they pass.
We have finished our first TDD cycle!
Time to commit the current set of changes
into a feature branch of our source control system!

Returning empty lists isn't that challenging.
Let's specify the behavior of our service when it contains some data.
To keep this listing short, some of the existing code is abbreviated with `...`.


```cucumber
Feature: Listing todo items
  ...

  Scenario: the user has no todo items
    ...

  Scenario: the user has some todo items
    Given a running "todo" service with the todo items:
      | TEXT                       | USER ID |
      | finish Exosphere tutorial  | 1       |
      | write my own Exosphere app | 1       |
    When receiving a "list" command with the payload:
      """
      {
        userId: 1
      }
      """
    Then it replies with "listing" and the payload:
      """
      [
        { text: 'finish Exosphere tutorial'  },
        { text: 'write my own Exosphere app' }
      ]
      """
```

What other scenarios could occur here?
The service could contain todo items of other users.
Let's add another scenario
and think about how our service should behave in that situation:

```cucumber
Feature: Listing todo items
  ...

  Scenario: the user has no todo items
    ...

  Scenario: the user has some todo items
    ...

  Scenario: many users have todo items
    Given a running "todo" service with the todo items:
      | TEXT                       | USER ID |
      | finish Exosphere tutorial  | 1       |
      | write my own Exosphere app | 1       |
      | cook rice                  | 2       |
      | go camping                 | 3       |
    When receiving a "list" command with the payload:
      """
      {
        userId: 1
      }
      """
    Then it replies with "listing" and the payload:
      """
      [
        { text: 'finish Exosphere tutorial'  },
        { text: 'write my own Exosphere app' }
      ]
      """
```


> TODO: hardcore TDD action here


Takeaway:
> Exosphere makes it very easy to work on services.
> They can be scaffolded in a fully working form,
> and be worked on completely by themselves via TDD.

Our todo-service works well enough
to be used in the application now.
Next we are going to integrate it with the web server.


<table>
  <tr>
    <td><a href="11_integration_into_web_server.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
