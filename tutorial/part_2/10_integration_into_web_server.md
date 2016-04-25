<table>
  <tr>
    <td><a href="09_todo_service.md">&lt;&lt; building the todo service</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="02_create_internal_service.md">creating an internal service &gt;&gt;</a></td>
  </tr>
</table>


# Integrating with the web server

<table>
  <tr>
    <td>
      <b><i>
        Status: beta - basics implemented, needs more hands-on testing
      </i></b>
    </td>
  </tr>
</table>

The next step is to integrate our fully functioning todo service into the web server.

__~/todo-app/web-server/server.js__

```javascript
const ExoRelay = require('exorelay');
const express = require('express');

const exoRelay = new ExoRelay({serviceName: process.env.SERVICE_NAME,
                               exocomPort: process.env.EXOCOM_PORT});
exoRelay.listen(process.env.EXORELAY_PORT);
exoRelay.on('error', (err) => { console.log(`Error: %{err}`); });
exoRelay.on('online', () => {

  const app = express();
  app.set('view engine', 'jade');

  app.get('/', (req, res) => {
    exoRelay.send('todos.list', (todos) => {
      res.render('index', {todos: todos});
    });
  });

  app.listen(3000, () => {
    console.log('Todo web server listening on port 3000');
  });

});
```

With everything in place,
let's test if the integration between web server and todo service works:

```
$ cd ~/todo
$ exo run
```

When you go to [localhost:3000](http://localhost:3000),
refresh the page,
and look at the terminal output of the exo runner,
you see that the web server now makes a request to the todo service
before it renders the home page:

```
...
exorun  application ready
exocom  web  --[ tweets.list ]->  tweets
exocom       (no payload)
tweets  listing tweets: 0 found
exocom  tweets  --[ tweets.listed ]->  web
exocom          { count: 0,
exocom            tweets: [] }
   web  GET / 200 308.709 ms - 886
   web  GET /assets/main.js 200 1.876 ms - 13019
   web  GET /favicon.ico 200 5.942 ms - 318
```

The Exo runner provides details about ongoing requests and communication
on the command line.
This is helpful for debugging issues.
You see exactly which service sends which message to which other one.


## Adding Todos

There are no todos in the database yet.
Let's make this possible, add a few entries, and show them on the homepage:

> * add input field and button to view
> * configure controller to POST to /todos
> * add /todos URL to web
> * implement "add" action in todo service
> * configure /todos controller to redirect to "/"

Once we restart the application,
refresh the browser,
we can add todo entries.
Yay, the application is starting to do something!


## Homework

* add a way to delete todo items to the application


<table>
  <tr>
    <td><a href="11_helper_apps.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
