<table>
  <tr>
    <td><a href="02_web_server.md">&lt;&lt; the web server service</a></td>
    <th>Exosphere Design Goals</th>
    <td><a href="02_create_internal_service.md">creating an internal service &gt;&gt;</a></td>
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
But for such a small micro-service a full web server stack
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

The Exosphere SDK provides frameworks called exoservices
to build micro-services out of simple handler functions,
comparable to [AWS Lambda](https://aws.amazon.com/lambda).

They are extremely minimalistic and only consist of:
* a configuration file that defines the service name and the messages it sends/receives
* a simple way to define handler functions for all incoming messages and send replies
* methods to send their own messages to other services.

We will use the
[Exoservice framework for Node.JS](https://github.com/originate/exoservice-js)
to build our todo service in the next chapter.


<table>
  <tr>
    <td><a href="09_todo_service.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
