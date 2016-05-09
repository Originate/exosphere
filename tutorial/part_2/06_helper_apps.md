<table>
  <tr>
    <td><a href="05_web_server.md">&lt;&lt; the web server service</td>
    <th>Helper Applications</th>
    <td><a href="07_communication.md">communication &gt;&gt;</a></td>
  </tr>
</table>


## Helper Applications

<table>
  <tr>
    <td>
      <b><i>
        Status: beta - basics implemented, waiting for feedback
      </i></b>
    </td>
  </tr>
</table>

In the last chapter,
Exosphere performed a number of commands
that aren't specified in the service configuration,
like installing the [NPM](https://www.npmjs.com) dependencies of the web server,
or starting it.
These activities are implemented as CLI scripts
in the service's `bin` directory.
This makes them available to both Exosphere as well as the developer while coding.
Just add `./bin` to your PATH environment variable,
and you can run them directly from a service directory.
Thanks to this convention,
one can work on a large variety of services without having to learn different command-line tools.
Here are the different commands that can live in the `bin` directory:

<table>
  <tr>
    <th>build</th>
    <td>builds the service, i.e. compiles the source code into an executable</td>
  </tr>
  <tr>
    <th>lint</th>
    <td>runs the linters</td>
  </tr>
  <tr>
    <th>run</th>
    <td>boots up the service</td>
  </tr>
  <tr>
    <th>setup</th>
    <td>installs all dependencies for the service, so that it is ready to be booted up</td>
  </tr>
  <tr>
    <th>spec</th>
    <td>runs the tests</td>
  </tr>
  <tr>
    <th>watch</th>
    <td>
      builds continuously in the background when files change
    </td>
  </tr>
</table>

All of this is just a convention.
You can also specify the commands to run in your service configuration file.

Takeaway:
> Exosphere provides a convention
> to provide actionable information about services
> as directly runnable command-line scripts.


<table>
  <tr>
    <td><a href="07_communication.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
