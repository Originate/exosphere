module.exports =

  before-all: (done) ->
    done!


  ping: (_, {reply}) ->
    reply 'pong'


  greet: (payload, {reply}) ->
    reply 'greeting', "Hello #{payload.name}"


  sender: (_payload, {send}) ->
    send 'greetings', 'from the sender service'

  'ping ponger': (_, {reply}) ->
    reply 'ping pong'

  'some salutation': (_, {reply}) ->
    reply 'this salutation', 'salutations'
