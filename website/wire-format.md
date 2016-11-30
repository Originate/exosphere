# ExoSphere Communication Wire Format

<table>
  <tr>
    <td>
      <b><i>
        status: deprecated - this will be replaced with ZeroMQ
      </i></b>
    </td>
  </tr>
</table>


The message protocol used in this prototype
is optimized for extreme simplicity.
This is to encourage experimentation
and implementations in many different languages.

Future versions of Exocom might use more performant communication protocols.


## Sending messages

To send an Exosphere message:
* make an HTTP __POST__ request to

  ```
  http://localhost:<EXOCOM_PORT>/send/<message name>
  ```

* `EXOCOM_PORT` is given to you via an environment variable with the same name
* provide this JSON structure in the body of your request:

  ```json
  {
    "sender": "<name of your service>",
    "id": "<the UUID you just made up for this request>",
    "payload": ["any", "data", "you", "want", "to", "send", "with", "the", "message"]
  }
  ```
* If the format of your request was correct,
  i.e. the sending service is authorized to send this message,
  and the request body is valid,
  ExoCom will respond with an HTTP __200__ code.
* If you are trying to send an unknown message, ExoCom responds with a __404__.
* If you are trying to send a message that exists in the application,
  but this service is not authorized to send it, ExoCom responds with a __403__.
* If your request was not correctly formatted, ExoCom responds with a 400.
* ExoCom bugs cause a 500 response
* Keep in mind that the response to your HTTP request only encodes the result
  of your attempt of _sending_ the message.
  It does not mean the message was understood and correctly processed by the
  services it will be sent to.
  Their replies to your message will arrive via separate messages sent to you.


## Receiving messages

To receive Exosphere messages:
* listen at the port given to you via the environment variable `EXORELAY_PORT`
* messages to you are sent to via __POST__ requests to

  ```
  /run/<message name>
  ```

* the body of the POST request is a JSON structure with this schema:

  ```json
  {
    "id": "<id of this incoming message>",
    "responseTo": "<id of the message for which this message is a reply>",
    "payload": {
      "foo": "bar"
    }
  }
  ```
* the `id` field contains the UUID of this incoming request
* if this message is a reply to a command you sent out earlier,
  the `responseTo` field contains the UUID of your command to which this is a reply
* the `payload` field contains data accompanying this message, in JSON format


## Communication libraries

You don't have to manually implement the low-level communication infrastructure described above.
The Exosphere SDK provides libraries called __ExoRelays__ that do this for you.
The following ExoRelay implementations are available:

* [ExoRelay-JS](https://github.com/Originate/exorelay-js): ExoRelay for Node.JS
