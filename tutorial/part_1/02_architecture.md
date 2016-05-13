<table>
  <tr>
    <td><a href="01_design_goals.md">&lt;&lt; design goals</a></td>
    <th>Exosphere architecture</th>
    <td><a href="03_installation.md">installation &gt;&gt;</a></td>
  </tr>
</table>


# Exosphere Architecture

Exophere consists of:

* An open-source __[developer SDK](https://github.com/Originate/exosphere-sdk)__
  that allows to build Exosphere applications on Windows, OS X, and Linux computers.
  It helps with scaffolding, testing, debugging, runing, and deploying Exosphere code bases.
  The SDK can also be used to run full Exosphere stacks
  on your own production infrastructure.

* __component marketplace:__
  Exosphere applications are collections of loosely coupled services.
  These services can come from a a variety of sources,
  including Originate's component marketplace,
  any Git repository,
  or any folder on your development machine.

* __cloud runtime:__
  The Exosphere Cloud Runtime is operated by Originate
  and provides professional, cost-efficient hosting and maintenance of Exosphere apps
  with professional support, additional security and analytics features,
  and availability guarantees.


## Levels

There are practical limits to the size of software projects.
We can only add so many layers of abstraction
on top of the foundation we build upon
before we exceed the available time or budget,
or our teams and code bases becomes so large and complex
that they collapse under their own weight.
This means the sophistication of our applications
is determined by the
level of the foundation that we build upon.
We are - and always will be -
dwarfs standing on the shoulders of giants<sup>1</sup>.

Exosphere has 3 main layers that build on top of each other.
Each layer itself consists of a generic foundation based on a popular open source technology,
on top of which is more opinionated, high-level support.

<img src="02_layers.png" width="395" height="401" alt="architecture layers">

From bottom to top:

<table>
  <tr>
    <th>1a</th>
    <td>
      At the basic level Exosphere is just a <b>cloud runtime for Docker images</b>.
      This is where you run legacy applications.
      An Exosphere application is made out of a number of code bases running together
      as groups of Docker images.
    </td>
  </tr>
  <tr>
    <th>1b</th>
    <td>
      On top of the Docker runtime is an opinionated <b>Heroku-like PaaS</b>
      that provides popular databases, cacheing services,
      and other resources to running code.
    </td>
  </tr>
  <tr>
    <th>2a</th>
    <td>
      A generic <b>messaging infrastructure</b> for asynchronous communication
      of services running inside Docker.
    </td>
  </tr>
  <tr>
    <th>2b</th>
    <td>
      Opinionated, high-level support to treat services as reusable <b>network objects</b>
      interacting with each other via message passing.
    </td>
  </tr>
  <tr>
    <th>3a</th>
    <td>
      A generic extension of the micro-service and message bus architecture
      to the front-end
      (web browsers, mobile applications) to leverage the same benefits there
      and provide a foundation for real-time applications.
    </td>
  </tr>
  <tr>
    <th>3b</th>
    <td>
      Opinionated, high-level support for quickly assembling
      applications from full-stack application parts
    </td>
  </tr>
</table>

You can build on top of any layer,
thereby choosing a particular balance
of efficiency through high-level support
vs freedom and flexibility to do things however you want to.

Takeaway:
> Exosphere is an opinionated and high-level platform for developing and operating cloud-native and mobile applications
> layered on top of generic and extensible infrastructure.


<table>
  <tr>
    <td><a href="03_installation.md"><b>&gt;&gt;</b></td>
  </tr>
</table>


<hr>

<ol>
  <li>
    Those "giants" are really just other dwarfs, standing on more dwarfs.
    Its dwarfs all the way down. There are no giants.
    You can participate in this like anybody else.
    Even the smallest contribution is huge
    if others can build on top of it!
  </li>
</ol>
