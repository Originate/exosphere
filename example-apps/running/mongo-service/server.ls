module.exports =

  'mongo.list': (_, {reply}) ->
    console.log "'mongo' service received message 'mongo.list'"
    reply 'mongo.listed' do
      * name: 'Jean-Luc Picard'
      * name: 'William Riker'
