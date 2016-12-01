<table>
  <tr>
    <td><a href="08_todo_service.md">&lt;&lt; todo service</a></td>
    <th>Integrating with the web server</th>
    <td><a href="../readme.md">back to part II &gt;&gt;</a></td>
  </tr>
</table>


# Integrating with the web server

The next step is to integrate our fully functioning todo service into the web server.
If you get lost, you can find the full application at this particular state [here](code_09).
First we update the home page to show all todo entries.

<a class="runMarkdown_createFileWithContent">
__web/app/controllers/index-controller.js__

```js
class IndexController {

  constructor({send}) {
    this.send = send
  }

  index(req, res) {
    this.send('todo.list', {}, (todos) => {
      res.render('index', {todos})
    })
  }

}

module.exports = IndexController
```
</a>

The only difference is that we now send out a `todo.list` message in the "index" action,
and then render the view with its result provided as the variable `todos`.
Here are the corresponding updates to the view:

<a class="runMarkdown_createFileWithContent">
__web/app/views/index.jade__

```jade
extends layout

block content

  h2 Exosphere Todos list
  p Your todos:
  ul
    each todo in todos
      li= todo.text

  h3 add a todo
  form(action="/todos" method="post")
    label text
    input(name="text")
    input(type="submit" value="add todo")
```
</a>

We loop over the `todos` variable provided by the controller
to render the list of todo entries,
and show a form to create a new todo entry below.

With this in place,
let's test if the integration between web server and todo service works:

<a class="runMarkdown_consoleWithDollarPrompt">
```
$ cd ~/todo
$ exo run
```
</a>

When you open [localhost:3000](http://localhost:3000) in your browser,
and look at the output of the exo runner in the terminal,
you see that the web server now makes a request to the todo service
before it renders the home page:

```
...
exorun  application ready
exocom  web  --[ todo.list ]->  todo
exocom       (no payload)
  todo  listing todos: 0 found
exocom  tweets  --[ todo.listing ]->  web
exocom          []
   web  GET / 200 308.709 ms - 886
...
```

The Exo runner provides details about ongoing requests and communication
on the command line.
This is helpful for debugging issues.
You see exactly which service sends which message to which other one.

In this case, the web server sends a `todo.list` message to the todo service
to request a list of all todo entries.
The todo service logs some activity on its own,
then sends back a `todo.listing` message to the web server,
saying that it doesn't have any todo entries at this point.


## Adding todo items

Let's fix that by adding the ability to create todo entries!
First, we need a controller to add todos via the web UI:

<a class="runMarkdown_createFileWithContent">
__web/app/controllers/todos-controller.js__

```js
class TodosController {

  constructor({send}) {
    this.send = send
  }

  create(req, res) {
    this.send('todo.create', req.body, () => {
      res.redirect('/')
    })
  }

}
module.exports = TodosController
```
</a>

It sends a `todo.create` message
with the content of the submitted HTML form
to the todo service,
then redirects to the home page.

We need to create a route for the new controller:

<a class="runMarkdown_createFileWithContent">
__web/app/routes.js__

```js
module.exports = ({GET, resources}) => {
  GET('/', { to: 'index#index' })
  resources('todos', { only: ['create', 'destroy'] })
}
```
</a>

We also need to tell the Exosphere framework
that the web service now sends and receives messages:

<a class="runMarkdown_createFileWithContent">
__web/service.yml__

```yaml
name: web
description: serves HTML UI for the test app

setup: npm install
startup:
  command: node app
  online-text: web server is running

messages:
  sends:
    - todo.create
    - todo.list
  receives:
    - todo.created
    - todo.listing
```
</a>

That's it!
Restart the web server by stopping it with Ctrl-C and starting it again.
Now we can add new todos via the web UI! Check it out: open localhost:3000 in your browser.
The console also provides good coverage
of the message traffic within our micro-service application.

Takeaway:
> An Exosphere application can be built up gradually.
> Exosphere's message bus provides helpful information
> to develop and debug information flow across services
> during that process.

This concludes the introduction to cloud-native backends written using the Exosphere framework.
Enjoy writing your own cloud-native micro service application backends!


## Homework

Our todo app only allows to add todos, but not delete them yet.
As homework, please add a way to delete todo items to the application.


<table>
  <tr>
    <td><a href="/tutorial/part_2/readme.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
