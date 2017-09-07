const {bootstrap} = require('exoservice')
const {MongoClient} = require('mongodb')
const N = require('nitroglycerin')

bootstrap({
  'before-all': function (done) {
    MongoClient.connect(getMongoAddress(), N(function() {
      console.log("MongoDB connected")
      done()
    }))
  }
})


function getMongoAddress() {
  return "mongodb://" + (process.env.MONGO || 'localhost:27017')
}
