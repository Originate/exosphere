<table>
  <tr>
    <td><a href="readme.md">&lt;&lt; part I overview</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="02_architecture.md">architecture &gt;&gt;</a></td>
  </tr>
</table>


# History and Design Goals

Originate specializes in building a broad variety of
mobile, web, and cloud-based software solutions
in high quality, at high velocity, for a fixed budget and schedule.
We also operate our applications at scale,
and maintain them for many years.
We have a proven methodology and best practices
to do this successfully:

<table>
  <tr>
    <th>modern</th>
    <td>
      The software industry is evolving at a breathtaking pace.
      Leverage new technologies to your advantage,
      they are often orders of magnitudes more productive.
    </td>
  </tr>
  <tr>
    <th>polyglot</th>
    <td>
      There isn't a single software stack that addresses all our needs.
      Mix and match them to combine their advantages.
    </td>
  </tr>
  <tr>
    <th>modular</th>
    <td>
      Complexity becomes the #1 bottleneck after a while.
      Break it up into much simpler and manageable pieces,
      and assemble them into loosely coupled architectures.
    </td>
  </tr>
  <tr>
    <th>reusable</th>
    <td>
      Build software components generic so that they become reusable in other projects.
    </td>
  </tr>
  <tr>
    <th>TDD</th>
    <td>
      Work against dedicated functional, performance, security, and quality goals.
      This ensures you reach them in the shortest possible time,
      and remain there as the application evolves.
    </td>
  </tr>
  <tr>
    <th>automated</th>
    <td>
      Automate as much boilerplate as possible:
      setup, testing, deployment, analytics, handling emergencies
    </td>
  </tr>
  <tr>
    <th>open source</th>
    <td>
      Leverage the combined expertise of the global developer ecosystem
      instead of trying to re-invent the wheel,
      share your contributions back
      for improved integration, functionality, performance, and security.
    </td>
  </tr>
  <tr>
    <th>agility</th>
    <td>
      Optimize your development processes around the humans in them.
      Cut red tape to get things done productively as a team.
    </td>
  </tr>
  <tr>
    <th>lean</th>
    <td>
      Focus on the essential.
      Start with a
      <a href="http://blog.codeclimate.com/blog/2014/03/20/kickstart-your-next-project-with-a-walking-skeleton">walking skeleton</a>
      of your application,
      iteratively add important functionality while always having a shippable and testable product.
    </td>
  </tr>
</table>

But the biggest bottleneck after a while is always **complexity**.
It accumulates over time,
and is the underlying reason
why the industry currently needs the top 10% of its developer talent
to run even average software development projects with a high chance for success.
Breaking up complexity into pieces that are manageable for most developers
is the biggest driver for Exosphere's design.

With all this in mind, and reinventing the wheel on too many projects,
we decided to build the framework we always wanted to have.
One that automates best practices,
and maximizes the efficiency of the human factor
by helping to:
* set up code and infrastructure
* build features in a modular, maintainable, and reusable way
  against explicit and automatically verified
  functional, performance, security, and quality specifications
* semi-automatically deploy into a variety of environments
  that can be spun up and maintained easily
* operate products at scale with built-in security and analytics
* keep the whole setup clean, maintainable, extensible, and fun to work on over time
  despite occasional drastic changes
* work as many small, effective teams
  including product owners, testers, developers, and other stakeholders

We ended up with a development framework that defines:
* a __development model__ to build complex applications
  as a set of loosely coupled, maintainable, reusable services
* a __collection of pre-existing services__ that allows to quickly plug together powerful solutions
  out of off-the-shelf components,
  with the ability to replace them with custom-built versions as the need arises
* a __devops toolchain__
  that automates the typical development and operate steps
  throughout the project lifecycle
* a __cloud runtime__
  that runs the applications at scale
  and provides built-in monitoring, security, and analytics

This development framework is built on top of modern open-source technologies
and cloud infrastructure,
but goes way beyond that.
So we called it _Exosphere_,
which is the highest layer of Earth's athmosphere,
way above the clouds.


## Non-goals

Exosphere is not ...

* a new web server or mobile application stack, test framework, or CI server.
  Use your own, in whatever language.
* a re-invention of existing open-source infrastructure technologies or billion-dollar startups.
  Instead, it builds higher levels of devops support on top of the most popular existing ones.


Takeaway:
> Exosphere is a framework for building and operating service-based cloud applications.
> It addresses many pain points around large-scale software development.


<table>
  <tr>
    <td><a href="02_architecture.md"><b>&gt;&gt;</b></td>
  </tr>
</table>
