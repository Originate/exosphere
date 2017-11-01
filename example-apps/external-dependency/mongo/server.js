const {bootstrap} = require('exoservice')
const {MongoClient} = require('mongodb')
const N = require('nitroglycerin')

bootstrap({
  beforeAll: function (done) {
    MongoClient.connect(getMongoAddress(), {autoReconnect: true, reconnectTries: 60, reconnectInterval: 2000}, N(function() {
      console.log("MongoDB connected")
      done()
    }))
  }
})

function getMongoAddress() {
  return "mongodb://" + (process.env.MONGO || 'localhost:27017')
}
