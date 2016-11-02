module.exports =

  'users.list': (_, {reply}) ->
    console.log "'users' service received message 'users.list'"
    reply 'users.listed' do
      * name: 'Jean-Luc Picard'
      * name: 'William Riker'
