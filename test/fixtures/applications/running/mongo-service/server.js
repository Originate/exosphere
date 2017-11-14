const {bootstrap} = require('exoservice')

bootstrap({
  beforeAll: function(done) {
    done()
  },
  'mongo.list': function(_, {reply}) {
    console.log("'mongo' service received message 'mongo.list'")
    reply('mongo.listed', ([
      {name: 'Jean-Luc Picard'},
      {name: 'William Riker'}
    ]))
  }
})
