# ExoCom Developer Guidelines

This monorepository contains the following subprojects:
* __exocom-server:__ the actual bus implementation
* __exocom-mock-*:__ mock implementations of the ExoCom server
                     for testing
* __exorelay-*:__ client SDKs for writing services that talk to ExoCom
* __exoservice-*:__ full-stack frameworks
                    for writing micro-services
                    as pure business logic,
                    i.e. as _lambda functions_


## Terminology

An Exocom application consists of a number of collaborating services.
For example, the _todo-application_ consists of a _user_ service for storing user accounts,
a _todo_ service for storing todo items,
an _html_ service for exposing the application's functionality via an HTML UI,
a _rest_ service for exposing it via a REST API,
and a _chat_ service for exposing it via a chatbot.

"User" and "todo" are the __role__ that the respective services play in the application.
A role can be played by different __service types__.
For example, you could use the _FacebookUserService_ service type
to play the role of the _user_ service for your application in production,
and you could use a _MockUserService_ to play the same role for local testing.
You can also use the same service type for different roles.
For example, you could use a generic _MongoService_ as a quick-and-dirty prototype
for both the _usersService_ and the _todoService_ role,
since they both store data.

To run the application, Exosphere creates an instance of the message bus,
plus an instance of the respective service for each role.
Or several instances for load balancing.
These instances are called __clients__ (because they are clients of the message bus).
Clients have unique names to identify them, for example "web #1", "web #2", etc.
Clients communicate via the central __message bus__,
played by the [Exocom-Server](exocom-server).

The __application configuration__ file (`application.yml`) defines
which roles the application contains, and links to the services that play each role
in the different environments (production, development, etc).
This link can be a link to a Docker image of the service
or its source code on Github or locally.

Each service type contains a __service configuration__ file (`service.yml`),
which defines the __message types__ this service type will send and receive.
A _message type_ specifies the title, payload structure,
and semantic meaning for messages.

Combining the roles from the _application configuration_
and the messages from the _service configurations_
creates the __routing__ information for the application.
It specifies which _role_ can send and receive which _message type_.
Based on that routing information,
Exocom knows to which service instances (clients) messages sent by other services should be forwarded.

To allow for reusable services, Exocom can perform __message translation__.
For example, the _MongoService_ above listens to `add entry` and `get entry` messages.
When used for the _usersService_ role, however,
it needs to respond to `add user` and `get user` messages.
The mapping from `add user` to `add entry` is specified in `application.yml`
and performed at runtime by Exocom.

In production, there can be different __protection levels__ for service roles.
__Public__ service roles are visible to the internet,
while __private__ roles are located in a separate subnet not visible to the outside world.
Exocom let's services of all protection levels communicate.


## Installation

* set up dev environment:
  * install [Morula](https://github.com/Originate/morula)
  * run `morula all bin/setup`

* run tests for all subprojects:

  ```
  $ morula all bin/spec
  ```

* run tests for changed subprojects (when on a feature branch):

  ```
  $ morula changed bin/spec
  ```


## Publish a new version

* log into Docker using the `originate` account: `docker login`
* run `bin/publish <patch|minor|major>`
* to use the new ExoCom version,
  publish a new [Exosphere](https://github.com/Originate/exosphere) version
