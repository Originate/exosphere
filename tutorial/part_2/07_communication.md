<table>
  <tr>
    <td><a href="06_helper_apps.md"><b>&lt;&lt;</b> helper applications</a></td>
    <th>Inter-Service Communication Architecture</th>
    <td><a href="08_communication_format.md">communication format <b>&gt;&gt;</b></a></td>
  </tr>
</table>


# Inter-Service Communication Architecture

<table>
  <tr>
    <td>
      <b>
      <i>
        status: beta - rough version implemented, needs more testing and solidifications
      </i>
      </b>
    </td>
  </tr>
</table>


Next we are going to add the ability to store todo lists to our application.
In a service-oriented architecture this will be implemented separately from the web server,
since storing todo items is a different responsibility than serving HTML pages to browsers.
Conceptually, a request to the homepage,
which is now supposed to contain a list of todo items,
would follow this rough workflow:

<table>
  <tr>
    <td width="400">
      <img alt="simple architecture for step 3" src="07_architecture_simple.png" width="386" height="355">
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
          The homepage is now supposed to contain a todo list.
          The <i>web server</i> service asks
          the <i>Todo service</i> for that list.
          <br>&nbsp;
        </li>
        <li>
          The Todo service loads the list from its database
          and replies to [2] with the todo list data.
          <br>&nbsp;
        </li>
        <li>
          The web server can now render the HTML for the home page
          and reply to the HTTP request made in [1].
        </li>
      </ol>
    </td>
  </tr>
</table>

Before we implement this,
let's take a look at how services interact with each other in Exosphere.


## Asynchronous Communication

Traditional web services are accessed with normal, direct network requests,
for example to a REST interface.
This is simple and fast,
and works well if there are only a few services.
As we scale to dozens or hundreds of them,
this synchronous communication model starts to become painful, though.
Each service:
* needs to figure out where all the other services are
  and keep this address list up to date
* handle edge cases like re-sending a message when the receiver doesn't respond,
  for example because it is currently being upgraded or restarted
* report sent messages to a logging server to support monitoring, debugging, and real-time analytics

Asynchronous communication solves many of these issues.
In Exosphere, communication is therefore decoupled in:
* space: The sender doesn't send a message to a particular address.
         It broadcasts it,
         and the messaging infrastructure routes it to the right receivers
         based on its content and metadata.

* time: A command isn't directly answered in the same network request.
        It is buffered until the receiver can get to processing it,
        and the receiver sends one (or multiple) replies as separate messages,
        at a later time.


## Exocom

Communication between services requires a lot of boilerplate activities.
Each service must:
* know how to talk to every other service's API in its own specific way,
  in regards to authentication, data formats, URL endpoints etc.
* enforce security, like verifying who the sender of a message was,
  whether it was allowed to send that message,
  and negotiate a way to encrypt messages along the way that works for both.
* monitor traffic patterns and raise alarms for suspicious changes
  that might for example indicate a hacker attack

If each service would have to do all this by itself,
they would all do it slightly differently and inconsistently.
And with only 20 services we already have 190 potential connections!
This quickly becomes inefficient and frustrating.

Such infrastructural functionality should be centralized and standardized
by the framework.
Exosphere does this by providing a generic messaging framework.
It corresponds to layer __2a__ in
[Exosphere's layer model](../part_1/02_architecture.md#levels).
Exosphere services talk to each other over a message bus called __exocom__,
which is short for <b>Exo</b>sphere <b>com</b>munication.
Here is how the workflow described above would be implemented in Exosphere:

<table>
  <tr>
    <td width="300">
      <img alt="architecture for step 3" src="07_architecture.png" width="293" />
    </td>
    <td>
      <ol>
        <li>
          The user browses to our homepage.
          To display that page,
          her browser makes an HTTP request for it.
          <br>&nbsp;
        </li>
        <li>
          The web server requests the list of todo items
          by sending a <code>todos.list</code> message to ExoCom
          <br>&nbsp;
        </li>
        <li>
          ExoCom sends this message to the Todo service,
          because it knows that this service can understand it.
          <br>&nbsp;
        </li>
        <li>
          The Todo service loads the todo entries from its database and
          replies to ExoCom with a <code>todos.listing</code> message.
          <br>&nbsp;
        </li>
        <li>
          ExoCom forwards this reply to the web server service.
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
don't need to know about each other anymore,
nor need to be able to talk directly to each other?
The just talk to ExoCom.
This greatly simplifies implementing the communication logic in them.

Only ExoCom needs to know the addresses of services,
and being a part of the framework that deploys services
it has direct access to this information.
ExoCom centralizes logic to verify
whether a service is allowed to send or receive a particular message.
The services can now be isolated from each other in separate network partitions
for security reasons,
ExoCom re-sends messages when services crash,
and monitors and analyzes the ongoing traffic, load patterns, response times, and failure rates,
and makes this data available to monitors and auto-scalers for services.
Exocom is also a great place to debug application logic as it flows across services,
since pretty much all application traffic goes through it.


## Exorelays

The Exosphere SDKs provides drivers
that encapsulate the logic
for talking to ExoCom.
They are called __ExoRelays__,
since they relay messages between services.
Each service contains exactly one Exorelay instance.

<img src="07_exorelays.png" width="714" height="220">

Exosphere provides ExoRelays for most popular languages.
It is easy to write additional Exorelays for your stack.


## Message Bus Types

There is a variety of exocom implementations,
each one specialized for a different set of requirements:
* __[ExoCom-dev](https://github.com/originate/exocom-dev):__
  A lightweight in-memory message bus implementation with very low latency,
  for local development and small to medium-sized production traffic.
* __ExoCom-prod:__
  A horizontally scalable production-grade bus with built-in persistence,
  optimized for throughput on large-scale deployments.
* __ExoCom-enterprise:__
  A production-grade bus for security-sensitive industries.
  It transmits message payloads encrypted.


Takeaway:
> Exosphere provides powerful communication infrastructure
> between the services of an application.

Next, we are going to look at the format of messages sent via exocom.


<table>
  <tr>
    <td><a href="08_communication_format.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
