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
This is your swiss army knife for working with Exosphere.
It provides commands<sup>1</sup> to:
* create applications
* add/remove services to/from applications
* add/remove end points to/from services
* spin up new environments in the Exosphere cloud
* test/deploy/monitor services and applications to environments in the cloud

To keep things as simple as possible,
this tutorial puts the application in a subdirectory of your home directory.
You can of course work in any other other directory on your machine.
Just replace `~` with your directory path in all code examples throughout this tutorial.
Open a terminal, and run these two commands:

```
cd ~
exo create app todo
```

This command generates a number of folders and files for us.
We will look at them in a minute.
First, why are we even using a code generator here?


<table>
  <tr>
    <td><a href="02_scaffolding.md"><b>&gt;&gt;</b></td>
  </tr>
</table>


<hr>

<sup>1</sup>
An overview of all exo commands is given [here](../../scaffolding.md).
No need to race ahead, though,
we will go through each one in this tutorial.
