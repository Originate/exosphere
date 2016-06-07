<table>
  <tr>
    <td><a href="readme.md"><b>&lt;&lt;</b> part II overview</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="02_scaffolding.md">scaffolding <b>&gt;&gt;</b></a></td>
  </tr>
</table>


# The "exo" Tool

<table>
  <tr>
    <td>
      <b><i>
      status: idea - not implemented yet, need feedback
      </i></b>
    </td>
  </tr>
</table>


We will use Exosphere's powerful command-line application -
the `exo` command -
to create a shell of our todo application.
This tool is your swiss army knife for working with Exosphere.
It provides commands to:
* create applications
* create new versions of existing applications
* add/remove services to/from applications
* add/remove end points to/from services
* spin up new environments in the Exosphere cloud
* test/deploy/monitor services and applications to environments in the cloud

An overview of all exo commands is given [here](../../scaffolding.md).
No need to race ahead, though,
we will go through each one in this tutorial.

To keep things as simple as possible,
this tutorial puts the application in a subdirectory of your home directory.
You can of course work in any other other directory on your machine.
Just replace `~` with your directory path in all code examples throughout this tutorial.
Open a terminal, and run these two commands:

```
cd ~
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
    <td>(press [Enter] to accept the default value of 0.0.1)</td>
  </tr>
</table>

The generator creates a folders `~/todo-application` for us.
We will look at it in a minute.
First, the takeaway of this chapter:

> Exosphere provides a powerful command-line application
> that allows to work on and operate Exosphere applications.

Now, why are we even using a code generator here?


<table>
  <tr>
    <td><a href="02_scaffolding.md"><b>&gt;&gt;</b></td>
  </tr>
</table>
