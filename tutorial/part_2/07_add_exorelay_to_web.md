<table>
  <tr>
    <td><a href="06_communication.md"><b>&lt;&lt;</b> communication</a></td>
    <th>Inter-Service Communication</th>
    <td><a href="08_scaffolding_services.md">scaffolding services <b>&gt;&gt;</b></a></td>
  </tr>
</table>


## Adding an Exorelay to the web server

Our web server is written in Nodes.JS,
so we use [ExoRelay-JS](https://github.com/originate/exorelay-js) here.
Let's add it to our web server:

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


<table>
  <tr>
    <td><a href="08_scaffolding_services.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
