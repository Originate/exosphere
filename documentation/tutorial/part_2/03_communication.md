<table>
  <tr>
    <td><a href="02_web_server.md">&lt;&lt; the web server service</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="02_create_internal_service.md">creating an internal service &gt;&gt;</a></td>
  </tr>
</table>


# Communication

<table>
  <tr>
    <td>
      status: beta - mostly implemented, needs robustness improvements
    </td>
  </tr>
</table>


Next we will build the service that stores the todo lists and entries.
Integrating it into our application will result in this architecture:

<table>
  <tr>
    <td width="280">
      <img alt="simple architecture for step 3" src="03_architecture_simple.png" width="265">
    </td>
    <td>
      <ol>
        <li>
          The user browses to our homepage.
          To display that page,
          her browser makes an HTML request for it.
          <br>&nbsp;
        </li>
        <li>
          The homepage is now contains the todo list.
          To render it, the web server service asks
          the Todo service for that list.
          <br>&nbsp;
        </li>
        <li>
          The Todo service replies to [2] with a list of todo entries.
          <br>&nbsp;
        </li>
        <li>
          With the todo entries available,
          the web server renders the HTML for the home page
          and replies to the HTTP request made in (1).
        </li>
      </ol>
    </td>
  </tr>
</table>


## ExoCom

Let's talk about how this communication is implemented in Exosphere.

Exosphere applications end up having a lot (dozens, hundreds) of services
that need to communicate with each other.
If all these services would try to talk directly to each other,
all of them would have to:

* figure out where all the other services are and keep this address list up to date
* know how to talk to every service's API in its own specific way
* enforce security, like verifying who the sender of a message was,
  whether it was allowed to send that message,
  and negotiate a way to encrypt messages along the way that works for both.
* monitor traffic patterns and raise alarms for suspicious ones
  that might indicate a hacker attack
* handle edge cases like re-sending a message when the receiver doesn't respond

If each service would do all this by themselves,
each one would do it slightly differently and inconsistently.
This would be wasteful, inefficient, and frustrating to work with.
Its better to centralize and standardize such infrastructure-like functionality.

Exosphere does via an infrastructure service called __ExoCom__.
It allows the different application services to communicate with each other.
Here is how our the real communication architecture including ExoCom looks like:

<table>
  <tr>
    <td width="280">
      <img alt="architecture for step 3" src="03_architecture.png" width="256" />
    </td>
    <td>
      <ol>
        <li>
          The user browses to our homepage.
          To display that page,
          her browser makes an HTTP request for it.</li>
          <br>&nbsp;
        <li>
          The web server requests the list of todo items
          by sending a `todos.list` message to ExoCom
          <br>&nbsp;
        </li>
        <li>
          ExoCom sends this message to the Todo service
          <br>&nbsp;
        </li>
        <li>
          The Todo service replies to ExoCom with the list of todo entries.
          <br>&nbsp;
        </li>
        <li>
          ExoCom sends this reply to the web server service.
          <br>&nbsp;
        </li>
        <li>
          The web server renders the HTML for the home page
          and replies to the HTTP request made in [1].
        </li>
      </ol>
    </td>
  </tr>
</table>

This looks more complex on the surface, but it is actually simpler,
especially when there are more than 2 services involved.
Notice how how both the web server and the Todo service
don't need to know about each other anymore?
Nor need to be able to talk directly to each other?
The just talk to ExoCom.
This greatly simplifies implementing the communication logic in them.
Only ExoCom needs to know the address of services.
It also makes communication secure:
ExoCom centralizes logic to verify
whether a service is allowed to send or receive a particular message.
The services can be in separate network partitions
to isolate them from each other
for security reasons,
ExoCom re-sends messages when services crash,
and monitors the ongoing traffic and load patterns.


## Adding communication to existing code bases

The Exosphere SDKs provides libraries for talking to ExoCom.
They are called __ExoRelays__,
since they relay messages between services.
ExoRelays exist for most popular languages.
In this tutorial we use [ExoRelay-JS](https://github.com/originate/exorelay-js)
for Node.JS code bases.

Let's add an ExoRelay to our web server:

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
  and create an instance of it called `exoRelay`
* we configure it with a number of environment variables,
  which are provided by Exosphere
* once we have connected our ExoRelay instance to ExoCom,
  we take our web server online as usual
* in the request handler for `/`, we send out a `todos.list` message
  and give it a callback function
* when the response arrives, it runs our callback,
  which renders the `index` view and gives it the `todos` variable as payload
