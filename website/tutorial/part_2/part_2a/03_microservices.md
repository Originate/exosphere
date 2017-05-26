<table>
  <tr>
    <td><a href="02_app_config.md"><b>&lt;&lt;</b> the application configuration </a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="04_html_server.md">the html server service <b>&gt;&gt;</b></a></td>
  </tr>
</table>


# Microservices

Complexity grows into the biggest bottleneck
on large software projects<a href="#1"><sup>1</sup></a>:
* Complexity accumulates over time, and is costly to get rid of.
* As code bases become large,
  they need extra patterns to structure and organize all the code.
  This makes them even more complex.<a href="#2"><sup>2</sup></a>
* The inevitable
  [technical drift](http://blog.codeclimate.com/blog/2013/12/19/are-you-experiencing-technical-drift/)
  causes larger and larger amounts of cleanup in your growing code base,
  to the point where you don't have the time to do it anymore,
  and have to live with these issues forever.

As a result,
most large software projects turn into complex headaches
after a few years.
Ramp-up time to learn all the architecture boilerplate increases to weeks and months.
This in turn makes it harder to change teams and
exacerbating the _mythical man-month_.
Only your best developers have the big picture
and know the right place to put things,
causing even faster quality degradation of your code base.
Large projects miss deadlines and go over budget,
if they achieve their functional and performance goals at all.
Because of this,
the software industry is the only one where _legacy_ is a bad word.

Micro-service architectures address this
by breaking up applications
into many independent code bases.
Each code base (service) has one responsibility
and is therefore small, simple, and easy to ramp up and work on, and test.
It is designed, developed, stored, tested, and deployed on its own.

This follows the way biology organizes complexity.
A long time ago nature found out that
trying to build 6-foot tall single-cellular organisms
that live for 80 years simply doesn't work.
Instead, we are collections of billions of specialized - but individual - cells.
This has numerous advantages:
* Each cell doesn't have to be perfect anymore.
  It is okay if they accidentally gunk up after a few days.
  Our body constantly recycles old cells and replaces them with new ones,
  and we don't even notice it.
* Homeostasis (auto-balancing) causes our bodies to
  optimize for the particular workload we have in our lives.
  Frequently used body parts become stronger (get more cells),
  less used parts get reduced
  so that we can spend our energy in areas more useful to us.

The same advantages apply to breaking up huge and complicated code bases
into many specialized micro-services:
* Clearer separation of concerns allows more teams to work better in parallel.
* Each service is simple, and can be worked on with little training
  by a smaller and more efficient team.
  This allows to throw more people later in the project at problems,
  and thereby eliminates the mythical man-month effect to some degree!
* Services can be more specialized for their workload,
  for example by using different languages or libraries best suited for their task.
  Each service can be tuned for its individual traffic pattern,
  instead of trying to tune one huge code or data base for several different
  (and possibly incompatible) traffic patterns.
* Code doesn't have to be perfect anymore.
  It is okay if each service only survives a few hours in production
  before memory leaks accumulate,
  and the overall application will still be robust and dependable.


Takeaway:
> Microservices eliminate the #1 problem of large-scale software development,
> accumulated complexity.

Next, let's add the first service to our application!

<table>
  <tr>
    <td><a href="04_html_server.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>


<hr>

<a name="1"><sup>1</sup></a>
On the microscopic level,
in particle physics and chemistry,
gravity is the least powerful of the [fundamental forces](https://en.wikipedia.org/wiki/Fundamental_interaction).
It is so weak that it is pretty much irrelevant inside and between molecules, atoms, and their nuclei.
The electromagnetic, strong, and weak nuclear forces are so much stronger
that they define these worlds.
But unlike the latter forces,
gravity accumulates above the microscopic level.
So, on the macroscopic (planetary or galactic) level,
it becomes the strongest and defining force for everything.

The same is true for complexity in code bases.
Methods and classes are often so simple that complexity inside them is not a particular big concern,
and it is often neglected in the name of "getting things done".
But when uncontained,
complexity will accumulate,
and as our systems grow,
become the by far largest problem we spend time on.

<a name="2"><sup>2</sup></a>
See [FizzBuzzEnterpriseEdition](https://github.com/EnterpriseQualityCoding/FizzBuzzEnterpriseEdition)
