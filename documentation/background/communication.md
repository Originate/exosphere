# Communication in Exosphere Applications

Communication in Exosphere is fully asynchronous and decoupled in
* _space:_ the sender does not need to know the address of the receiver.
  It simply shouts out the message and trusts that the right receivers will listen to it.
* _time:_ the receiver doesn't have to be online or ready to receive when the message is sent.
  The message waits in a queue until the receiver is ready to process it.

This means messages are fire-and-forget.
You send off a message, and if there is a response, you will receive it
as a separate message.


## Terminology

Services send each other __messages__.
To send a message, a service makes a __request__ to
[ExoCom](https://github.com/Originate/exocom-dev).
The __response__ for this request only indicates whether
the message was correctly sent,
not whether it was executed successfully.
ExoCom then makes other requests to the receiving services
in order to deliver the message to them.

There are two types of messages:
* __commands__ are normal messages sent to another service
  in order to make it do something.
  An example is the web server sending the users service the _"user.create"_
  command to create a new user account.

* __replies__ are messages sent in response to a previously received command.
  They contain the __outcome__ of the command that triggered them.
  In our example, the users service would reply to a received _"user.create"_ command
  with a _"user.create.ok"_ message containing the created user account.


## Message Format

Messages have the following structure:

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
    <td>version of the message</td>
    <td>1.2</td>
  </tr>
  <tr>
    <th>payload</th>
    <td>JSON</td>
    <td>data payload for the message</td>
    <td>{"name": "Jean-Luc Picard"}</td>
  </tr>
  <tr>
    <th>id</th>
    <td>UUID</td>
    <td>guid of the message</td>
    <td>7294be20-e034-11e5-bbcc-1d18f38c4e43</td>
  </tr>
  <tr>
    <th>requestId</th>
    <td>UUID</td>
    <td>guid of the over-arching end-user request</td>
    <td>23d09070-e098-11e5-91a5-a7f02faca148</td>
  </tr>
  <tr>
    <th>response-to</th>
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

The structure of the payload can be defined using
[JSON schema](http://json-schema.org).


## Communication Protocol

The low-level communication protocol for ExoCom is described [here](wire-format.md).



## Outcomes

Service calls are higher-level, usually crossing functional boundaries
within your application, performing dedicated units of work.
And going across the network makes them more expensive
than in-process function calls
in a monolithic application.
Hence they typically (should) have more complex APIs than function calls.

* replies to commands often include the state changes caused by the command,
  to avoid having to do round trip to the service to query the new state
* commands often have more than one outcome.
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

The boundaries between "success" and "no success",
between "user error" and "application error" are blurry here.
Is a pending transaction a success? It doesn't feel like an error,
but the money also hasn't been transferred as requested.
Is a "transaction limit exceeded" outcome a user error?
Could/should the app have warned the user before it tried the transaction?
Because of this, its better to model outcomes in the application's domain language,
and configure monitoring and alerts in the application's analytics system.


## Streaming Replies

Services can send more than one reply to a command.
In this case we have a _message stream_ that encodes a _streaming reply_.

One use case for streaming replies is monitoring of longer-running commands.
Let's say we have a "blob" service that stores large files.
Copying a large file takes time, so we want to display a progress bar.
To make this easy, the "blog" service sends streaming replies to the "file.copy" command.
Here is how this looks like using [ExoRelay](https://github.com/Originate/exorelay-js):

```livescript
exoRelay.send 'file.copy', from: 'large.csv', to: 'backup.csv', (payload, {outcome}) ->
  switch outcome
    | 'file.copy.in-progress'  =>  console.log "still copying, #{payload.percent}% done"
    | 'file.copy.done'         =>  console.log 'file copy finished!'
```

Another use case is streaming responses, where a larger result is sent in a series of chunks:

```livescript
exoRelay.send 'file.read', path: 'large.csv', (payload, {outcome}) ->
  switch outcome
    | 'file.read.chunk'  =>  result += payload
    | 'file.read.done'   =>  console.log "finished reading #{payload.megabytes} MB!"
```


## Subscription Model

Services declare the messages they send and receive in their configuration file (config.yml).
When deployed, Exosphere reads this configuration and sets up the right
subscriptions and channels for each service.


Exosphere guarantees delivery of every message to at least one instance of
each subscribed service type.

Messages can be delivered to multiple instances,
for example if instances crash while processing the message, or fail to confirm
having received and processed it.
