const {bootstrap} = require('exoservice')
const {MongoClient} = require('mongodb')

bootstrap({
  beforeAll: function (done) {
    MongoClient.connect(getMongoAddress(), {autoReconnect: true}, function(error) {
      if (error) {
        console.error(error)
      } else {
        console.log("MongoDB connected")
        done()
      }
    })
  }
})

function getMongoAddress() {
  return "mongodb://" + (process.env.MONGO || 'localhost:27017')
}
