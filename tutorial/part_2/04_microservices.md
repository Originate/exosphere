<table>
  <tr>
    <td><a href="03_app_config.md"><b>&lt;&lt;</b> the application configuration </a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="05_web_server.md">microservices <b>&gt;&gt;</b></a></td>
  </tr>
</table>


# Microservices

<table>
  <tr>
    <td>
      <b><i>
      status: beta - feature complete, needs fine-tuning
      </i></b>
    </td>
  </tr>
</table>

Over time,
complexity always grows into the biggest bottleneck
on large software projects:
* Complexity accumulates over time, and is very hard to get rid of.
* At some point the code base becomes so large
  that it needs a lot of extra patterns to organize all the millions of lines of code.
* Ramp-up time to learn all this architecture boilerplate increases to weeks and months,
  making it harder to change teams and
  exacerbating the _mythical man-month_.
* The inevitable
  [technical drift](http://blog.codeclimate.com/blog/2013/12/19/are-you-experiencing-technical-drift/)
  causes larger and larger amounts of cleanup in your growing code base,
  to the point where you don't have the time to do it anymore
  and have to live with these issues forever.
* Only a few developers except your best people have the big picture
  and know the right place to put things,
  causing even faster quality degradation of your code base.

As a result,
most large software projects turn into complex headaches
after a few years,
go over budget if they end successfully at all,
and the software business is the only industry where _legacy_ is a bad word.

Exosphere addresses this by
breaking up applications
into many independent code bases.
Each code base (service) has one responsibility
and is therefore small, simple, and easy to work on and test.
It is designed, developed, stored, tested, and deployed on its own.

This follows the way biology organizes complexity.
A long time ago nature found out that
trying to build 6-feet tall single-cellular organisms
that live for 80 years doesn't work.
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
  so that we can spend their energy in areas more useful to us.

The same advantages apply to breaking up a huge complicated code base
into many specialized micro-services.
* Clearer separation of concerns allows more teams to work better in parallel.
* Each service is simple, and can be worked on with little training
  by a smaller and more efficient team<sup>1</sup>.
* Services can be more specialized for their workload,
  for example by using the language or libraries best suited for their task.
  Each service can be tuned for its individual traffic pattern,
  instead of trying to tune one huge code or data base for 50 separate patterns.
* Code doesn't have to be perfect anymore.
  It is okay if each service only lasts a few hours before memory leaks accumulate,
  and the overall application is still robust and dependable.


Next, let's add the first service to our application!

<table>
  <tr>
    <td><a href="05_web_server.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
