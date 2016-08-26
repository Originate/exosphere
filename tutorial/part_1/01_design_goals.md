<table>
  <tr>
    <td><a href="readme.md">&lt;&lt; part I overview</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="02_architecture.md">architecture &gt;&gt;</a></td>
  </tr>
</table>


# History and Design Goals

Originate specializes in building a broad variety of
AI-native, mobile, web, and cloud-based software solutions
in high quality, at high velocity, and hitting scoped budget and timelines with a very high probability.
We work with different languages, frameworks, and cloud runtimes.
We also operate our applications at scale,
and evolve them for many years.

We built a framework to integrate and automate the many best practices
we have developed to do this reliably,
and provide missing features we always wanted to have available
to do our job well and have fun along the way.
We named it Exosphere, because it goes way beyond the clouds.

Exosphere is a combination of developer SDK, cloud runtime, and service bazaar
for building, operating, and evolving
service-oriented, AI-native application suites.

An Exosphere application starts out as a set of functional, performance, security, and quality metrics
expressed as automatically executable documentation, specifications, and tests.
These metrics drive the development and evolution of the application.

functional components of applications
as sets of independent, reusable, loosely coupled services,
and provides a way for these services to interact
by exchanging high-level, application specific messages.


Exosphere provides:
* infrastructure to define
  and work against them
* automated code and infrastructure setup and deployment


## Goals


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

The Exosphere framework is not ...

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
