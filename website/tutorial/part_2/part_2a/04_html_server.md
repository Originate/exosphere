<table>
  <tr>
    <td><a href="03_microservices.md"><b>&lt;&lt;</b> microservices</a></td>
    <th>The Web Server Service</th>
    <td><a href="05_communication.md">communication <b>&gt;&gt;</b></a></td>
  </tr>
</table>

# The HTML Server Service

Let's build the first service for our application: the HTML server! If we were
building our Todo app as a traditional monolith, this would be the only code
base and perform all of the application's functionality: receiving requests,
storing todo items, user accounts, and API tokens in the database, configuring
and using the search engine, tracking user sessions and permissions, rendering
HTML for the browser, JSON for the REST API, sending out emails, etc.

In Exosphere's microservice world, each code base has only one responsibility.
The html server's job is to interact with the user via an HTML UI. Most of the
things mentioned above are not a direct part of this responsibility, and are
therefore implemented outside of the html server, as separate services. Because
of this much narrower set of responsibilities, the html server is a lot smaller
and simpler than it would be in a traditional monolithic application.

Since our html server is so simple, we'll build it using
[ExpressJS](http://expressjs.com).

Exosphere provides a template for building ExpressJS html servers. So, before we
add the service, we will add the HTML server service template and then we will
add the new service.

<a class="runMarkdown_consoleWithInputFromTable">

```
cd ~/todo-app
git clone https://github.com/originate/exosphere-htmlserver-express .exosphere/service_templates/htmlserver-express
exo add service
```

Again, the generator asks for the information it needs interactively. Please
enter:

<table>
  <tr>
    <th>prompt</th>
    <th>text you enter</th>
  </tr>
  <tr>
    <td>Template</td>
    <td>1</td>
  </tr>
  <tr>
    <td>Service Role</td>
    <td>web-service</td>
  </tr>
  <tr>
    <td>App Name</td>
    <td>html-server</td>
  </tr>
  <tr>
    <td>Description</td>
    <td>serves the HTML UI of the Todo app</td>
  </tr>
  <tr>
    <td>Type</td>
    <td>public</td>
  </tr>
  <tr>
    <td>Author</td>
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
  web-service:
    location: ./web-service
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

The html service is located in a subdirectory of the application, in
`~/todo-app/web-service/`. This makes sense because it is an integral part of
our application, and doesn't make sense outside of it.

Most of the files in this folder are just a normal
[ExpressJS](http://expressjs.com) application, plus some extra tools like
linters.

Let's check out how the service looks like internally.

### The service configuration file

This file tells the Exosphere framework everything it needs to know about this
service.

**~/todo-app/web-service/service.yml**

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
  * the **online-text** section contains the string to look for in the terminal
    output to determine when the service has successfully started up. The
    Exosphere runtime only sends traffic to fully available instances.
* The **messages** section lists all the messages that this service will send
  and receive. The Exosphere runtime needs this information in order to
  automatically subscribe the service to these messages. Currently our
  application doesn't contain any other services that could be communicated
  with, so this section is empty for now. We'll add some commands here soon!

The other files in this directory are just a normal
[ExpressJS](http://expressjs.com) application.

## Booting up the application

With all files in place, the Exosphere CLI has all the information to set up our
application. Let's have Exosphere run our application for us:

```
exo run
```

The output should look something like:

```
>>> docker-compose --file run_development.yml up --build
Creating network "todoapp_default" with the default driver
Pulling exocom0.26.1 (originate/exocom:0.26.1)...
0.26.1: Pulling from originate/exocom
add3ddb21ede: Pull complete
7b0e69015745: Pull complete
Digest: sha256:56e176f25e21732ccad671375fba8177a04a3fc7ea8bc1decf9e715c6af961c6
Status: Downloaded newer image for originate/exocom:0.26.1
Building web-service
Step 1/6 : FROM node
latest: Pulling from library/node
85b1f47fba49: Pull complete
ba6bd283713a: Pull complete
817c8cd48a09: Pull complete
47cc0ed96dc3: Pull complete
8888adcbd08b: Pull complete
6f2de60646b9: Pull complete
9dd205971dc0: Pull complete
5859715a4691: Pull complete
Digest: sha256:7c9099e0f68242387d7755eaa54c287e16cedd3cca423444ca773794f5f1e423
Status: Downloaded newer image for node:latest
 ---> c1d02ac1d9b4
Step 2/6 : ENV NODE_ENV "development"
 ---> Running in 61ad323a3101
 ---> 0d4ade2ed6a6
Removing intermediate container 61ad323a3101
Step 3/6 : COPY package.json .
 ---> 382f908e332a
Step 4/6 : COPY yarn.lock .
 ---> 2b2387775f86
Step 5/6 : RUN yarn
 ---> Running in e91437d97867
yarn install v1.3.2
[1/4] Resolving packages...
[2/4] Fetching packages...
info fsevents@1.1.2: The platform "linux" is incompatible with this module.
info "fsevents@1.1.2" is an optional dependency and failed compatibility check. Excluding it from installation.
[3/4] Linking dependencies...
[4/4] Building fresh packages...
Done in 10.40s.
 ---> 714bc67a7dc4
Removing intermediate container e91437d97867
Step 6/6 : COPY . .
 ---> 5a6ff5c8cff8
Successfully built 5a6ff5c8cff8
Successfully tagged todoapp_web-service:latest
Creating exocom0.26.1 ...
Creating exocom0.26.1 ... done
Creating web-service ...
Creating web-service ... done
Attaching to exocom0.26.1, web-service
exocom0.26.1    | ExoCom online at port 80
exocom0.26.1    | 'web-service' registered
web-service     | ExoRelay for 'web-service' online at port 80
web-service     | HTML server online at port 3000
web-service     | HTML server is running
```

The Exosphere framework itself is written as a bunch of loosely coupled
services. We see a number of them in action here:

* **exocom** is the command that runs Exosphere applications. It starts the
  other services.
* **web-service** is our html server service. We can see that exorun starts it,
  and recognizes right after the output `Todo app running at port 3000` that our
  html server is online.
* The Exosphere runtime also starts a service called **exocom**. This is the
  messaging system for communication between services. More about it later.

Finally, exorun tells us that the application is now fully started and ready to
be used. Open a browser and navigate to
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
