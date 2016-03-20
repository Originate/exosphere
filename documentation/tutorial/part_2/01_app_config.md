<table>
  <tr>
    <td><a href="readme.md">&lt;&lt; part II overview</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="02_create_internal_service.md">creating an internal service &gt;&gt;</a></td>
  </tr>
</table>


# Application Configuration

In a microservice environment,
applications consist of many different code bases.
Each code base is stored in its own repository,
is worked on by its own team,
and deployed on its own schedule.

In return, each code base is very small and simple, and has one responsibility.
This setup allows to break up complexity and parallelize the work of teams better.
Since each service is so simple, it is also possible to work on them with very
little training.

An Exosphere application is not much more than
a configuration file that defines all the
code bases that make up the application.

Here is the configuration file of our Todo application

```yml
name: Todo application
description: Allows to store notes
version: '0.0.1'

services:
  web:
    location: ./web-server
```

It defines the name, a description, and a version of the application.
It also defines an initial service:
The "web" service provides the web UI for the application.
For each service it lists the place where to find it.
In this case the services are in subdirectories of the application.
This location could also point to a Git repository.

Later builds of this application could also contain an "api" service
as well as "user-auth" and "session" services to allow users to log in.
