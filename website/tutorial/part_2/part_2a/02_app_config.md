<table>
  <tr>
    <td><a href="01_exo_tool.md"><b>&lt;&lt;</b> the exo tool</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="03_microservices.md">microservices <b>&gt;&gt;</b></a></td>
  </tr>
</table>


# The Application shell

Since all the action within an application happens in the services,
an Exosphere application itself is not much more than
a configuration file that defines all the
services that make up the application.
We call this _composite applications_.

Here is the configuration file created for our Todo application:

<a class="tutorialRunner_verifyFileContent">

__todo-app/application.yml__
```yml
name: todo-app
description: An example Exosphere application
version: 0.0.1

development:
  dependencies:
    - name: exocom
      version: 0.26.1

services:
```

</a>

It defines the name, a description, and a version of the application.
It also defines a section for listing all the services of the application and
any development dependencies. In our case, this includes `exocom`.
Since we just created the application, the services section will be empty.
We will add some in just a minute.

First, the takeaway:
> An Exosphere application is nothing but some configuration data
> and a list of services that make up the application's functionality and data storage.

Lets talk about why our Exosphere application doesn't contain any code
in the next chapter!

<table>
  <tr>
    <td><a href="03_microservices.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
