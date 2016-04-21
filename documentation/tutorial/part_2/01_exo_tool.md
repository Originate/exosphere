<table>
  <tr>
    <td><a href="readme.md">&lt;&lt; part II overview</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="02_web_server.md">the web server service &gt;&gt;</a></td>
  </tr>
</table>


# The "exo" Tool

<table>
  <tr>
    <td>
      <b><i>
      status: idea - not implemented yet
      </i></b>
    </td>
  </tr>
</table>


We will use Exosphere's powerful command-line application -
the `exo` command -
to create a shell of our todo application.
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

Code generation (aka code scaffolding) and other forms of automation
play a much bigger role in microservice environments
than they do for monoliths like for example [Ruby on Rails](http://rubyonrails.org).
In monoliths, they are useful in only a few places,
and besides that are mostly used to impress beginners
how fast one can generate simple CRUD functionality
that is almost never used in the generated form.
That's because once a monolithic code base is set up,
adding a new feature often means to simply add a few files in a bunch of directories,
and that's often easier done by hand.

In the micro-service world,
adding a new feature often means setting up one or several new code bases,
each one with its own:
* source code shell
* configuration files
* test framework
* documentation
* code repository (the Git repo to contain this code base)
* CI server setup
* automatic deployment based on Git Flow or a comparable strategy
* integration into the main application and other services

Doing all this manually easily takes a few hours.
Without automation,
developers would be discouraged from setting up new services.
They would try to cram new functionality into an existing service
and then "clean this up later when we have more time".
That's an anti-pattern that you should avoid at all cost<sup>1</sup>.
Our code base would drift back into monolith land
where code bases are massive and complicated,
except that we wouldn't even have the support of a framework for monoliths.
In a micro-service environment,
intelligent automation is critical for maintaining high development velocity and code quality.
And its better doable than in monoliths,
since code in services is so much simpler.

The `exo` command-line tool is your swiss army knife for working with Exosphere.
It provides commands<sup>2</sup> to:
* create applications
* add/remove services to/from applications
* add/remove end points to/from services
* spin up new environments in the Exosphere cloud
* test/deploy/monitor services and applications to environments in the cloud

In the next chapter, we'll look at what our command above has created.
[>>](02_app_config.md)


<hr>

<sup>1</sup> Its a trap! Don't ever do this when writing production code.
There will never be a time when you will have nothing to do
and will be paid to go back and clean up existing code that already works.
Developer time is expensive and always in short supply.
In the future, there will be lots of other supert important things to do.
The only way to have a clean code base is to build it clean right from the start.
There is still a need to refactor your code when you experience
[technical drift](http://blog.codeclimate.com/blog/2013/12/19/are-you-experiencing-technical-drift/),
but try to avoid acculumating [technical debt](https://en.wikipedia.org/wiki/Technical_debt)
at all cost,
or your code base will turn into a mess sooner than you think,
and you'll be the one who has to live with it.

<sup>2</sup>
An overview of all exo commands is given [here](../../scaffolding.md).
No need to race ahead, though,
we will go through each one in this tutorial.
