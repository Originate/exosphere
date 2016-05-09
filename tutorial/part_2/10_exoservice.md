<table>
  <tr>
    <td><a href="09_add_exorelay_to_web.md">&lt;&lt; add exorelay to the web server</a></td>
    <th>Exoservices</th>
    <td><a href="11_todo_service.md">building the todo service&gt;&gt;</a></td>
  </tr>
</table>


## Exoservices

<table>
  <tr>
    <td>
      <b><i>
        Status: beta - basics implemented, needs more hands-on testing
      </i></b>
    </td>
  </tr>
</table>

To build a service,
we could whip together another quick [ExpressJS](https://expressjs.com) stack
and use an [ExoRelay](https://github.com/Originate/exorelay-js) inside it,
as we did in the [web server service](02_web_server.md).
But for a small micro-service a full web server stack
based on the MVC paradigm would be overkill:
* we don't need a sophisticated _routing_ layer,
  since micro-services don't serve a large variety of URLs:
  mostly, they just implement a single REST-like endpoint
* we don't need a sophisticated _model_ layer,
  since a service typically only deals with one or very few model types
* we don't need _views_, since our service returns simple JSON data
  that can be automatically serialized
* since we don't have models and views,
  we don't really need _controllers_ either

All we need are simple handler functions
that are called for incoming messages.
They do a little bit of database work and send the outcome back.
The Exosphere SDK provides frameworks called Exoservices
that allow to build such services in the most radical and minimalistic way.
Exoservices only consist of:
* a configuration file that defines the service name and the messages it sends/receives
* a way to define handler functions for all incoming messages
* methods to send replies as well as their own messages to other services

Takeaway:
> Exosphere allows to build services that consist almost completely of business logic,
> with boilerplate reduced to almost zero.

In the next chapter we will use the
[Exoservice framework for Node.JS](https://github.com/originate/exoservice-js)
to build our todo service.


<table>
  <tr>
    <td><a href="11_todo_service.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
