const {bootstrap} = require('exoservice')

bootstrap({

  beforeAll: (done) => {
    done()
  },

  // Replies to the "ping" command
  ping: (_, {reply}) => {
    reply('pong')
  }

})
