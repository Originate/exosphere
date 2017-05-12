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
  if process.env.MONGO
    return "mongodb://#{process.env.MONGO}"
  return "mongodb://localhost:27017"
