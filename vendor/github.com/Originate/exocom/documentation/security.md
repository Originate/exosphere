# ExoCom Security Model

ExoCom centralizes cross-cutting concerns like security and logging.
This frees developers from re-implementing those aspects redundantly in every service.
Instead, they can rely on those aspects being provided by the runtime.

## How it works

- services transmit messages to be sent to ExoCom
- ExoCom sends an "can send?" message including metadata about the message
  to its security service instance over a private channel
- the security service can be written in any language, and runs in Docker
- the security service determines based on it's own business logic
  whether this message should be allowed to be sent
- the security service instance replies with either a
  "can send" or "cannot send" message to the request from the ExoCom instance
- depending on the answer, ExoCom distributes the message to the subscribing service
  instances or logs a security incident

[[ image of the architecture ]]
