<img src="documentation/logo.png" width="862" height="111" alt="logo">

[![Build Status](https://travis-ci.org/Originate/exosphere.svg?branch=master)](https://travis-ci.org/Originate/exosphere)

_yak shaver for coders_

Exosphere is an infrastructure framework for composite code bases
that automates typical repetitive activities
that software developers do:

- [clone](exo-clone/features) a composite application
  consisting of code located in a variety of repositories
  onto your machine
- [add](exo-add/features) more services to your application -
  either a fully functioning empty service including tests from templates,
  or an existing third-party service from an external repository or DockerHub image
- [run](exo-run/features) the application on your machine.
  Exosphere boots up all services and their dependencies
  within Docker images, so no installation
  of programming languages or runtimes is necessary.
  It is built on top of [Docker Compose](https://docs.docker.com/compose/)
  and adds missing features on top of it.
  For example waiting until dependencies are fully booted up
  before starting the application services.
  Or restarting individual Docker images
  as their content gets modified by the developer
- [test](exo-test/features) all the services of the application,
  as well as the application as a whole.
  This includes functional, performance<sup>&#42;</sup>,
  reliability (chaosmonkey style)<sup>&#42;</sup>,
  and security tests<sup>&#42;</sup>.
<!--- update all the third-party application parts to their latest version --->
- [deploy](exo-deploy) your application to public or private cloud environments
  including AWS and others<sup>&#42;</sup>.

Exosphere combines industrial-strength open-source technologies
into an opinionated bundle that makes them work well together:
* [Git](https://git-scm.com) for source code management
* [Docker](https://www.docker.com) for containerization
* [Boilr](https://github.com/tmrts/boilr) for code generation
* [Terraform](https://www.terraform.io) for deployment


## Why

Breaking up monoliths into composite applications
(which can consist of dozens to hundreds of indvidiual code bases)
can improve complexity management, decoupling, and reusability,
but also makes a lot of things more complicated and repetitive.
For example, now you have to:
- download many individual code bases onto your machine
- prepare them, for example by installing modules or compiling them
- install several runtimes or compilers for the different code bases,
  sometimes in conflicting versions
  that are not trivial to install and operate in parallel on the same machine.
- install and configure dependencies like databases, routers, caches
- deploying now means
  testing, dockerizing, uploading, deploying, and configuring
  dozens/hundreds of individual code and data bases.

There are good tools for many of these tasks,
but they don't know about each other and don't work together.
Many of them also lack features.
As a result, projects often duct-tape a variety of tools together,
each one in their own ways.
Everybody on the team has to install, configure, and learn to use
complicated tools before being able to
get the application running and make changes to it.

Exosphere makes it easy to work on any code base,
without having to install, configure, master, and use
a whole array of devops tools
and learn how exactly they are used on each project.


### Benefits for tech leads

As a tech lead, you can make your knowledge and expertise available
in the form of a tool that you configure,
so that the rest of your team can follow your way
without you having to be around and supervise everything yourself.
You can focus on the next big question
instead of being busy keeping the project running.


### Benefits for developers

As a developer, you can get started working on the project within minutes,
and everything works.
Instead of having to read extensive project documentation that explains
what exact set of ops and infrastructure tools
and other boilerplate
that are not part of your business logic,
but are merely there
to help deal with complexity.
are used on this particular project,
in what particular way,
and having to install and learn to use them,


### Benefits for Ops

Exosphere makes Ops easier in several ways:
- Exosphere provides infrastructure that enforces dedicated
  functional, performance, reliability, and security goals for code,
  against which developers have to build their code against
  right from the start.
  This means the product will run better and more reliable in production,
  because these things are no longer afterthoughts.
- issues encountered in production, for example a performance or security regression,
  can be modeled as new restrictions and goals which are enforced to the developers
  by their toolchain. This means they will be addressed, and keep being addressed,
  from now on.
- the code you operate comes with a lot of devops tool-chain attached to it,
  which makes it easier to set up CI and CD
- Exosphere has built-in support to set up and deploy to various environments,
  for example development, qa, staging, production


### Business benefits

Saving your engineers from the need to deal with custom, duct-taped together
low-level ops tools provides a number of benefits from a business perspective:
- the ramp-up time for new engineers on the project is radically reduced,
  sometimes from weeks to minutes.
- you can utilize a wider variety of skill levels on the project.
  For example, with Exosphere it becomes so trivial to run even complex applications
  that even non-technical people like designers or product managers
  can run the application on their machine.
- your tech lead is less busy keeping the project in the air,
  and can focus on more strategic questions


## Get Started
* [learn more](website/config_files) about the configuration files
* [install the SDK](website/tutorial/part_1/03_installation.md)
* [download and run an example application](website/example-apps.md)
* build your own application by following the [tutorial](website/tutorial)


## Reference
* the [scaffolding commands](website/scaffolding.md)
* [platform developer documentation](website/developers/developers.md)


## Further Reading
* related projects: [LeverOS](https://github.com/leveros/leveros)


<hr>

<sup>&#42;</sup>
coming soon
