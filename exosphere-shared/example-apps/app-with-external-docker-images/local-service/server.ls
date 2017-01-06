module.exports =

  before-all: (done) ->
    done!


  # replies to the "ping" command
  ping: (_, {reply}) ->
    reply 'pong'
