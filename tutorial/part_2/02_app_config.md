<table>
  <tr>
    <td><a href="01_exo_tool.md">&lt;&lt; the exo tool</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="03_web_server.md">the web server service &gt;&gt;</a></td>
  </tr>
</table>


# The Application shell

<table>
  <tr>
    <td>
      <b><i>
      status: alpha - some parts are implemented, needs robustness improvements
      </i></b>
    </td>
  </tr>
</table>

In Exosphere's microservice world,
applications are broken up into many individual code bases.
Each code base (service) has one responsibility
and is therefore small, simple, and easy to work on and test.
It is stored on its own (often in its own repository),
and is tested and deployed by itself.

Breaking up a code base like this
prevents build-up of large and complex bodies of code (which are hard to understand)
and allows several teams to work better in parallel.
And since services have a much simpler structure than a monolithic code base,
it is possible to work on them with much less training and ramp-up time
than would be required to wrap one's head around a massive monolithic code code base
and all the complexity, abstraction, and patterns that are necessary to manage it<sup>1</sup>.


## The application configuration file

Since all the action within an application happens in the services,
an Exosphere application itself is not much more than
a configuration file that defines all the
services that make up the application.

Here is the configuration file created for our Todo application:

```yml
name: Todo application
description: Allows to store notes
version: '0.0.1'

services:
```

It defines the name, a description, and a version of the application.
It also defines a section for listing all the services of the application.
It is empty right now, since our application doesn't contain any services yet.


<hr>

<sup>1</sup>
In the grand scheme of things
a microservice application of course still contains this complexity.
It is pushed out of sight of the (many) developers who work on basic services
as the building blocks of an application.
It happens in business logic a level or two above them.
With micro-services,
these two aspects (building basic functionality and orchestrating it into an application)
are better separated than in a monolith,
hence architect-level expertise is only needed in a few places,
not everywhere,
and it is easier to put things in the right place
and thereby keep everything well maintained.
