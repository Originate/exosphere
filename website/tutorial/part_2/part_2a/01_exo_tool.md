<table>
  <tr>
    <td><a href="readme.md"><b>&lt;&lt;</b> part II overview</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="02_app_config.md">The application shell <b>&gt;&gt;</b></a></td>
  </tr>
</table>


# The "exo" Tool

We will use Exosphere's powerful command-line application -
the `exo` command -
to create a shell of our todo application.
This tool is your swiss army knife for working with the Exosphere framework.
It provides commands to:
* create applications
* add/remove services to/from applications
* add/remove end points to/from services
* bump versions of services and applications
* spin up new environments in the Exosphere cloud
* test/deploy/monitor services and applications to environments in the cloud

An overview of all exo commands is given [here](../../../scaffolding.md).
No need to race ahead, though,
we will go through each command in this tutorial.

Okay, let's get started!
Open a terminal, and run these two commands:

<a class="tutorialRunner_runConsoleCommand">

```
exo create application
```

The command asks for all necessary information interactively.
Please enter:

<table>
  <tr>
    <th>prompt</th>
    <th>text you enter</th>
  </tr>
  <tr>
    <td>Name of the application to create</td>
    <td>todo-app</td>
  </tr>
  <tr>
    <td>Description</td>
    <td>An example Exosphere application</td>
  </tr>
  <tr>
    <td>Initial version</td>
    <td>(Press enter and accept the default)</td>
  </tr>
    <tr>
    <td>ExoCom Version</td>
    <td>(Press enter and accept the default)</td>
  </tr>
</table>

</a>

The generator creates a folder `~/todo-app` for us.
We will look at it in a minute.
First, the takeaway of this chapter:

> The Exosphere framework provides a powerful command-line application
> that allows you to work on and operate Exosphere applications.

Now, let's take a closer look at the generated code.


<table>
  <tr>
    <td><a href="02_app_config.md"><b>&gt;&gt;</b></td>
  </tr>
</table>
