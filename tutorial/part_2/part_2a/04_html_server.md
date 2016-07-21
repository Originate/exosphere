<table>
  <tr>
    <td><a href="03_microservices.md"><b>&lt;&lt;</b> microservices</a></td>
    <th>The Web Server Service</th>
    <td><a href="05_communication.md">communication <b>&gt;&gt;</b></a></td>
  </tr>
</table>


# The HTML Server Service

Let's build the first service for our application:
the HTML server!
If we would be building our Todo app as a traditional monolith,
this would be the only code base
and perform all of the application's functionality:
receiving requests,
storing todo items, user accounts, and API tokens in the database,
configuring and using the search engine,
tracking user sessions and permissions,
rendering HTML for the browser, JSON for the REST API,
sending out emails,
etc.

In Exosphere's microservice world,
each code base has only one responsibility.
The html server's job is to interact with the user via an HTML UI.
Most of the things mentioned above are not a direct part of this responsibility,
and are therefore implemented outside of the html server,
as separate services.
Because of this much narrower set of responsibilities,
the html server is a lot smaller and simpler
than it would be in a traditional monolithic application.

Since our html server is so simple,
we'll build it using [ExpressJS](http://expressjs.com).
Exosphere provides a template for building ExpressJS html servers.
Let's use it:

```
cd ~/todo-app
exo add service
```

Again, the generator asks for the information it needs interactively.
Please enter:

<table>
  <tr>
    <th>prompt</th>
    <th>text you enter</th>
  </tr>
  <tr>
    <td>Name of the service to add</td>
    <td>html-server</td>
  </tr>
  <tr>
    <td>Description</td>
    <td>serves the HTML UI of the Todo app</td>
  </tr>
  <tr>
    <td>type</td>
    <td>htmlserver-express-es6</td>
  </tr>
  <tr>
    <td>Initial version</td>
    <td>(press [Enter] to accept the default value of 0.0.1)</td>
  </tr>
</table>

Now we see the service registered in our application configuration file:
__application.yml__

```yml
name: Todo application
description: An example Exosphere application
version: '0.0.1'

services:
  html-server:
    location: ./html-server
```

Here is the current architecture of our application:

<table>
  <tr>
    <td width="280">
      <img alt="architecture for step 2" src="04_architecture.png" width="258">
    </td>
    <td>
      <ol>
        <li>
          The user browses to our homepage.
          In order to show that page, her web browser requests the HTML for it.
        </li>
        <li>
          This request goes to our <i>html server service</i>.
          It replies with the HTML for the page.
        </li>
      </ol>
    </td>
  </tr>
</table>



## The html service folder

The html service is located in a subdirectory of the application,
in `~/todo-app/html-server/`.
This makes sense because it is an integral part of our application,
and doesn't make sense outside of it.

Most of the files in this folder
are just a normal [ExpressJS](http://expressjs.com) application,
plus some extra tools like linters.

Let's check out how the service looks like internally.


### The service configuration file

This file tells the Exosphere framework everything it needs to know about this service.

__~/todo-app/html-server/service.yml__

```yml
name: Todo html server
description: serves the html UI of the Todo app

setup: npm install --loglevel error --depth 0
startup:
  command: node app
  online-text: html server listening on port

messages:
  sends:
  receives:
```

This file tells the Exosphere framework about this service.
* The service __name__ and a __description__
* The __setup__ section defines the commands to make the service runnable on the machine,
  i.e. install its dependencies, compile it, etc.
* The __startup__ section defines how to boot up the service.
  * the __command__ section contains the command to run
  * the __online-text__ section contains the string to look for in the terminal output
    to determine when the service has successfully started up.
    The Exosphere runtime only sends traffic to fully available instances.
* The __messages__ section lists all the messages that this service will send and receive.
  The Exosphere runtime needs this information
  in order to automatically subscribe the service to these messages.
  Currently our application doesn't contain any other services
  that could be communicated with,
  so this section is empty for now.
  We'll add some commands here soon!


The other files in this directory are just a normal
[ExpressJS](http://expressjs.com)
application.


## Setting up the service

With all files in place,
the Exosphere CLI has all the information to set up our application.
Let's check that the overall configuration is correct,
and have Exosphere set up the service for us:

```
$ exo setup
```

We see how it uses Node's package management system (NPM)
to download and install
the external ExpressJS and Jade modules for us,
so that the service is ready to run.
The output should look something like:

```
Setting up Todo application 0.0.1

  exo-setup  starting setup of 'html-server'
html-server  /Users/kevin/exosphere-sdk/tutorial/part_2/code_02/todo-app/html-server
├── express@4.13.4
└── jade@1.11.0
html-server  PROCESS ENDED
html-server  EXIT CODE: 0
  exo-setup  setup of 'html-server' finished
  exo-setup  setup complete
```


## Booting up the application

To test that everything works, let's check that the application boots up:

```
$ exo run
```

You should see output like:

```
Running Todo application 0.0.1

 exocom      online at port 8000
html-server  Server running at port 3000
 exorun      'html-server' is running using exorelay port 8001
 exorun      all services online
 exocom      received routing setup
 exorun      application ready
```

The Exosphere framework itself is written as a bunch of loosely coupled services.
We see a number of them in action here:
* __exorun__ is the command that runs Exosphere applications.
  It starts the other services.
* __html-server__ is our html server service.
  We can see that exorun starts it,
  and recognizes right after the output `Todo app running at port 3000`
  that our html server is online.
* The Exosphere runtime also starts a service called __exocom__.
  This is the messaging system
  for communication between services.
  More about it later.

Finally, exorun tells us that the application is now fully started
and ready to be used.
Open a browser and navigate to [http://localhost:3000](http://localhost:3000).
We got a running microservice-based web site!

Takeaway:
> The web server in a microservice application is much simpler than in a monolith,
> because it only focuses on interacting with the user via HTML.

Next, let's look at how services communicate with each other!

<table>
  <tr>
    <td><a href="05_communication.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
