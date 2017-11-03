const {bootstrap} = require('exoservice')
const {MongoClient} = require('mongodb')

const connectWithRetry = function(retryDelay, done) {
  MongoClient.connect(getMongoAddress(), function(error) {
    if (error) {
      console.log(`Could not connect, retrying in ${retryDelay} milliseconds...`, error.toString())
      setTimeout(() => connectWithRetry(retryDelay * 2, done), retryDelay)
    } else {
      console.log("MongoDB connected")
      done()
    }
  })
};

bootstrap({
  beforeAll: function (done) {
    connectWithRetry(1000, done)
  }
})

function getMongoAddress() {
  return "mongodb://" + (process.env.MONGO || 'localhost:27017')
}
