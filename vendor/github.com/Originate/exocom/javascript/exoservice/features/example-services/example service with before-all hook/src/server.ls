hooks-ran = []


module.exports =

  before-all: (done) ->
    hooks-ran.push 'before-all'
    done!


  'which-hooks-ran': (_, {reply}) ->
    reply 'these-hooks-ran', hooks-ran
