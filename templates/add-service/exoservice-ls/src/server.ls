module.exports =

  before-all: (done) ->
    # TODO: add asynchronous init code here, or delete the whole block
    done!


  # replies to the "ping" command
  ping: (_, {reply}) ->
    reply 'pong'
