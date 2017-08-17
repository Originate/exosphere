require! {
  'mongodb' : {MongoClient}
  'nitroglycerin' : N
}

module.exports =

  before-all: (done) ->
    MongoClient.connect get-mongo-address!, N (mongo-db) ->
      console.log("MongoDB connected")
      done!

function get-mongo-address
  return "mongodb://#{process.env.MONGO or 'localhost:27017'}"
