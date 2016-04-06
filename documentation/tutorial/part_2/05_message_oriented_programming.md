<table>
  <tr>
    <td><a href="03_communication.md">&lt;&lt; communication</a></td>
    <th>Message-Oriented Programming</th>
    <td><a href="05_todo_service.md">building the Todo service&gt;&gt;</a></td>
  </tr>
</table>


# Message-oriented Communication

One (nice) side effect of having this application-wide communication architecture
is that communication becomes __message oriented__.
Services don't make random network transmissions to each other,
they send each other well-formatted, high-level __messages__.

Sending messages is asynchronous and fire-and-forget:
You send them to ExoCom and that's it.
No need to verify whether they have been received,
retrying them if something went wrong processing them,
or manually encrypt them to secure them from being snooped on.
Exosphere takes care of all of that for you
by providing a reliable, secure messaging layer.
It forwards your message to the services that are subscribed to it,
guarantees that only one instance of each service receives it,
retries sending it a number of times before giving up,
and transmits any responses back to you.


## Message-oriented Programming

Messages don't correspond to simple function calls in object-oriented programming.
Rather, they represent high-level service calls
across the functional boundaries of your code
that drive the control flow of your application.
To support this as much as possible,
they are built as statically typed, shared, immutable data structures.

