<table>
  <tr>
    <td><a href="02_scaffolding.md"><b>&lt;&lt;</b> scaffolding</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="04_microservices.md">microservices <b>&gt;&gt;</b></a></td>
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


Since all the action within an application happens in the services,
an Exosphere application itself is not much more than
a configuration file that defines all the
services that make up the application.

Here is the configuration file created for our Todo application:

```yml
name: Todo application
description: Allows to store notes
version: 0.0.1

services:
```

It defines the name, a description, and a version of the application.
It also defines a section for listing all the services of the application.
It is empty right now, since our application doesn't contain any services yet.
We will add some in just a minute.

First, the takeaway:
> An Exosphere application is nothing but some configuration data
> and a list of services that make up the application's functionality and data storage.

Lets talk about why our Exosphere application doesn't contain any code
in the next chapter!

<table>
  <tr>
    <td><a href="04_microservices.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>

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
