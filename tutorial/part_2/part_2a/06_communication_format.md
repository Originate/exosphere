<table>
  <tr>
    <td><a href="05_communication.md"><b>&lt;&lt;</b> communication</a></td>
    <th>Inter-Service Communication Format</th>
    <td><a href="07_exoservices.md">Exoservices <b>&gt;&gt;</b></a></td>
  </tr>
</table>


# Inter-Service Communication Format

In this chapter we look at the data that flows through
the communication architecture
we discussed in the [last chapter](06_communication.md).
How are messages sent over exocom structured?


## Terminology

First, let's clarify the vocabulary used,
since we deal with three types of requests now:

1. Our backend responds to __activities__ from its users.
  Activities are network calls
  coming in from the internet
  to one of the services exposed to the public,
  like our web or API server.
  Examples for activities are
  somebody visiting the home page of our application
  using a web browser,
  or a cron job making a call to our public API.

2. Typically, several services within our application work together
  to handle an activity.
  To do this, they exchange a number of __messages__.
  There are two types of messages:

  * __commands__ are messages sent to another service
    in order to make it do something.
    Typically, commands are in imperative ("create").
    An example is the web server sending the todo service a `todo.list`
    command to list all todo items.

  * __replies__ are messages sent in response to a previously received command.
    They reference the command that caused them,
    and contain the __outcome__ of what they did to satisfy the command.
    Replies are verbs,
    either in present continuous tense ("creating")
    to indicate that an activity is occurring right now,
    or in past tense ("created")
    to indicate that an activity has occurred and is finished now.

3. To send a message,
  a service makes some form of network __request__ to
  ExoCom.
  The __response__ for this network request only indicates
  whether the message was correctly sent,
  i.e. whether its payload was correctly formatted,
  not whether it was executed successfully by its receivers.
  Exocom then performs additional _requests_
  to the services that are supposed to receive the message
  in order to deliver the message to them.


## Outcomes

Commands, i.e. service calls, are higher-level operations
than simple function calls.
They cross functional boundaries within the application,
i.e. go from one part of the application (the web server) to a different part (the todo service).
They typically perform higher-level and complex work,
and can therefore have a larger variety of possible outcomes.

For example, a message encoding the command
_"transfer $100 from the checking account to the savings account"_,
sent to an account service, could produce any of these outcomes:

<table>
  <tr>
    <th>transferred</th>
    <td>the money was transferred</td>
  </tr>
  <tr>
    <th>pending</th>
    <td>the transfer was initiated, but is pending a third-party approval</td>
  </tr>
  <tr>
    <th>transaction limit exceeded</th>
    <td>the account doesn't allow that much money to be transferred at once</td>
  </tr>
  <tr>
    <th>daily limit exceeded</th>
    <td>the daily transaction limit was exceeded</td>
  </tr>
  <tr>
    <th>insufficient funds</th>
    <td>there isn't enough money in the checking account</td>
  </tr>
  <tr>
    <th>unknown account</th>
    <td>one of the given accounts was not found</td>
  </tr>
  <tr>
    <th>unauthorized</th>
    <td>the currently logged in user does not have privileges to make this transfer</td>
  </tr>
  <tr>
    <th>internal error</th>
    <td>an internal error occurred in the accounting service</td>
  </tr>
</table>

The boundaries between _success_ and _no success_,
between _user error_ and _application error_ become blurry here.
Is a pending transaction a success? It doesn't feel like an error,
but the money also hasn't been transferred as requested.
Is a "transaction limit exceeded" outcome a user error?
Should the computer have warned the user
before even allowing a transaction that it later determines will fail?
Many of these questions can only be decided by the application, not the account service.
Rather than having the account service try to determine what is an error and what is success,
its better to simply return the outcomes in correct domain language,
as shown in this table,
and then let the application interpret the outcome.


## Streaming Replies

Services can send more than one _reply_.
In this case we have a _message stream_ that encodes a _streaming reply_.

One use case are streaming responses,
where a larger result is sent in a series of chunks:
As an example, here is how a large file is read in chunks using
[ExoRelay.JS](https://github.com/Originate/exorelay-js):

```javascript
exoRelay.send('file.read', {path: 'large.csv'}, (payload, {outcome}) => {
  switch(outcome) {

    case 'file.read-chunk':
      result += payload
      break

    case 'file.read-done':
      console.log(`finished reading ${payload.megabytes} MB!`)
      break
  }
}
```

Another use case for streaming replies is making longer-running commands responsive.
Let's say we have a "blob" service that stores large files.
Copying a larger file might take more than a few seconds,
so we want to display a progress bar to the user.
To make this easy,
the "blog" service sends streaming replies to the "file.copy" command.

```javascript
exoRelay.send('file.copy', {from: 'large.csv', to: 'backup.csv'}, (payload, {outcome}) => {
  switch(outcome) {

    case 'file.copy.in-progress':
      console.log(`copying, ${payload.percent}% done`)
      break

    case 'file.copy.done':
      console.log('file copy finished!')
      break

  }
}
```

Note:
These are quick solutions for simple use cases.
Don't pump too much data over the message bus.
Exocom a control system, not a big data transfer system.
When reading very large files,
use a dedicated data transfer protocol,
and just transmit the address to read from over the bus.
If a file copy takes a longer time,
or other services want to display the file copy operation as well,
build your own dedicated transaction model
instead of relying on replies sent to a particular service instance.


## Message Format

Finally, let's look at how messages actually look like.
Each message contains this metadata in its header:

<table>
  <tr>
    <th></th>
    <th>type</th>
    <th>description</th>
    <th>example</th>
  </tr>
  <tr>
    <th>name</th>
    <td>string</td>
    <td>type of the message sent</td>
    <td>"user.create"</td>
  </tr>
  <tr>
    <th>version</th>
    <td>number</td>
    <td>version of the API talked to</td>
    <td>1.2</td>
  </tr>
  <tr>
    <th>id</th>
    <td>UUID</td>
    <td>guid of the message</td>
    <td>7294be20-e034-11e5-bbcc-1d18f38c4e43</td>
  </tr>
  <tr>
    <th>timestamp</th>
    <td>unix time in nanoseconds</td>
    <td>time when this message was sent (determined by exocom)</td>
    <td>987234987234987</td>
  </tr>
  <tr>
    <th>activityId</th>
    <td>UUID</td>
    <td>guid of the event that caused this message (user activity or cron job)</td>
    <td>23d09070-e098-11e5-91a5-a7f02faca148</td>
  </tr>
  <tr>
    <th>responseTo</th>
    <td>UUID</td>
    <td>id of the command replied to here</td>
    <td>89934740-e034-11e5-bbcc-1d18f38c4e43</td>
  </tr>
  <tr>
    <th>sender</th>
    <td>string</td>
    <td>name of the service that sends the message</td>
    <td>"web"</td>
  </tr>
</table>

The payload of the message can be provided in any way the applications wants to,
for example as JSON, [MsgPack](http://msgpack.org), using
[Protocol Buffers](https://developers.google.com/protocol-buffers), or
[Thrift](https://thrift.apache.org).
You can transmit textual, binary, or any other data in it.


## Type checking

You are encouraged to use a technology that codifies and verifies the structure
of exchanged messages, like protocol buffers.
No matter whether you implement your services statically or dynamically typed,
the control and data flow between different parts of your application should be verified.
Especially when the different services talking to each other are maintained by different teams.

Takeaway:
> Exosphere messages are designed to be high-level calls between subsystems of an application.
> They have a well defined format,
> can have multiple types and numbers of replies,
> include metadata,
> and allow any types of payload.


Next we'll implement communication in our todo application.


<table>
  <tr>
    <td><a href="07_exoservices.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
