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
If we were
building our Todo app as a traditional monolith, 
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
So, before we add the service, 
we will add the HTML server service template and then we will add the new service.

<a class="runMarkdown_consoleWithInputFromTable">

```
cd ~/todo-app
git clone https://github.com/originate/exosphere-htmlserver-express .exosphere/service_templates/htmlserver-express
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
    <td>template</td>
    <td>1</td>
  </tr>
  <tr>
    <td>role</td>
    <td>html-server</td>
  </tr>
  <tr>
    <td>name</td>
    <td>html-server</td>
  </tr>
  <tr>
    <td>description</td>
    <td>serves the HTML UI of the Todo app</td>
  </tr>
  <tr>
    <td>type</td>
    <td>public</td>
  </tr>
  <tr>
    <td>author</td>
    <td></td>
  </tr>
</table>

</a>

Now we see the service registered in our application configuration file:

<a class="runMarkdown_verifyFileContent">
__application.yml__

```yml
name: todo-app
description: An example Exosphere application
version: 0.0.1
development:
  dependencies:
  - name: exocom
    version: 0.26.1
services:
  html-server:
    location: ./html-server
```

</a>

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
type: public
description: serves the HTML UI of the Todo app
author:

# defines how to boot up the service
startup:
  # the string to look for in the terminal output
  # to determine when the service is fully started
  online-text: HTML server is running

# the messages that this service will send and receive
messages:
  sends:
  receives:

# other services this service needs to run,
# e.g. databases
dependencies:

docker:
  ports:

development:
  port: 3000
  scripts:
    # the command to boot up the service
    run: node ./index.js
```

This file tells the Exosphere framework about this service.

* The service **type**, **description**, and **author**
* The **startup** section defines how to boot up the service.
  * The **online-text** section contains the string to look for in the terminal output 
    to determine when the service has successfully started up. 
    The Exosphere runtime only sends traffic to fully available instances.
* The **messages** section lists all the messages that this service will send and receive. 
  The Exosphere runtime needs this information 
  in order to automatically subscribe the service to these messages. 
  Currently our application doesn't contain any other services 
  that could be communicated with, 
  so this section is empty for now. 
  We'll add some commands here soon!

The other files in this directory are just a normal
[ExpressJS](http://expressjs.com) application.

## Booting up the application

With all files in place, 
the Exosphere CLI has all the information to set up our application. 
Let's have Exosphere run our application for us:

```
exo run
```

The output should look something like:

```
>>> docker-compose --file run_development.yml up --build
Building html-server
Creating exocom0.26.1 ... done
Creating html-server ... done
Attaching to exocom0.26.1, html-server
exocom0.26.1    | ExoCom online at port 80
exocom0.26.1    | 'html-server' registered
html-server     | ExoRelay for 'html-server' online at port 80
html-server     | HTML server online at port 3000
html-server     | HTML server is running
```

We see a couple services get started here:

* **exocom** is the messaging system for communication between services. 
  More about it later.
* **html-server** is our html server service. 
  We can see that exorun starts it,
  and recognizes right after the output `Todo app running at port 3000` 
  that our html server is online.

Finally, exosphere tells us that the application is now fully started and ready
to be used. Open a browser and navigate to
[http://localhost:3000](http://localhost:3000). We got a running
microservice-based web site!

Takeaway:
> The web server in a microservice application is much simpler than in a
> monolith, because it only focuses on interacting with the user via HTML.

Next, let's look at how services communicate with each other!

<table>
  <tr>
    <td><a href="05_communication.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
