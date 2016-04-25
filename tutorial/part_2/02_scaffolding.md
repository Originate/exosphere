<table>
  <tr>
    <td><a href="01_exo_tool.md"><b>&lt;&lt;</b> the exo tool</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="03_app_config.md">the application shell <b>&gt;&gt;</b></a></td>
  </tr>
</table>


# Scaffolding

<table>
  <tr>
    <td>
      <b><i>
      status: idea - not implemented yet, need feedback
      </i></b>
    </td>
  </tr>
</table>


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
intelligent automation of typical development tasks
is critical for maintaining high development velocity and code quality.
And its easier doable than in monoliths,
since code in services is so much simpler.

Takeaway:
> Scaffolding allows to build applications efficiently in safe, incremental steps.


In the next chapter, we'll look at what our command above has created.


<table>
  <tr>
    <td><a href="03_app_config.md"><b>&gt;&gt;</b></td>
  </tr>
</table>


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
