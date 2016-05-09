<table>
  <tr>
    <td><a href="08_communication_format.md"><b>&lt;&lt;</b> communication format</a></td>
    <th>Inter-Service Communication</th>
    <td><a href="10_exoservice.md">Exoservices <b>&gt;&gt;</b></a></td>
  </tr>
</table>


## Adding an Exorelay to the web server

<table>
  <tr>
    <td>
      <b><i>
      Status: beta - basics implemented, needs more hands-on testing
      </i></b>
    </td>
  </tr>
</table>

Let's add an ExoRelay to our web server.
It is written in Nodes.JS,
so we use [ExoRelay-JS](https://github.com/originate/exorelay-js) here.
Let's add it to our code base:

```
$ cd ~/todo-app/web-server
$ npm install --save exorelay
```

Now use the ExoRelay in the web server code:

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

What we do differently now:
* we import the `ExoRelay` class from the `exorelay` NPM module
  and create an instance called `exoRelay`
* we configure it with a number of environment variables,
  which are provided by Exosphere
* once we have connected our ExoRelay instance to ExoCom,
  we take our web server online as usual
* in the request handler for the homepage,
  we send out a `todos.list` message
  and give it a callback function
* when the response to the `todos.list` command arrives,
  it runs our callback,
  which renders the `index` view and gives it the `todos` variable as payload


Takeaway:
> Exorelays can be added to any code base
> to enable it to communicate with other Exosphere services.

Next, we are going to build the Todo service.
Before we start,
let's discuss how we are going to do it!

<table>
  <tr>
    <td><a href="10_exoservice.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
