<table>
  <tr>
    <td><a href="10_integration_into_web_server.md">&lt;&lt; integration with the web server></td>
    <th>Service CLI</th>
    <td><a href="12_message_oriented_communication.md">message-oriented programming &gt;&gt;</a></td>
  </tr>
</table>


## service CLI

<table>
  <tr>
    <td>
      <b><i>
        Status: beta - basics implemented, needs more hands-on testing
      </i></b>
    </td>
  </tr>
</table>

In the last chapter, we used a command called `spec`.
It is part of a standardized set of CLI commands
that each service within Exosphere should define in its `bin` directory.
Since services in Exosphere can be written in any language and use any tool set,
a developer would have to learn a whole array of tools to install, boot, and test
the different services they use.
Thanks to this convention, you can use these commands across every service:

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


<table>
  <tr>
    <td><a href="09_todo_service.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
