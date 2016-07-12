require! {
  'mongodb' : {MongoClient, ObjectID}
  'nitroglycerin' : N
}
env = require('get-env')('test')


collection = null

module.exports =

  before-all: (done) ->
    mongo-db-name = "exosphere-_____serviceName_____-#{env}"
    MongoClient.connect "mongodb://localhost:27017/#{mongo-db-name}", N (mongo-db) ->
      collection := mongo-db.collection '_____modelName_____s'
      console.log "MongoDB '#{mongo-db-name}' connected"
      done!


  # Creates a new _____modelName_____ object with the given data
  '_____modelName_____.create': (_____modelName_____-data, {reply}) ->
    collection.insert-one _____modelName_____-data, (err, result) ->
      | err  =>
          console.log "Error creating _____modelName_____: #{err}"
          reply '_____modelName_____.not-created', error: err
      | _  =>
          console.log "creating _____modelName_____"
          reply '_____modelName_____.created', mongo-to-id(result.ops[0])


  '_____modelName_____.create-many': (_____modelName_____s, {reply}) ->
    collection.insert _____modelName_____s, (err, result) ->
      | err  =>  return reply '_____modelName_____.not-created-many', error: err
      reply '_____modelName_____.created-many', count: result.inserted-count


  '_____modelName_____.read': (query, {reply}) ->
    try
      mongo-query = id-to-mongo query
    catch
      console.log "the given query (#{query}) contains an invalid id"
      return reply '_____modelName_____.not-found', query
    collection.find(mongo-query).to-array N (_____modelName_____s) ->
      switch _____modelName_____s.length
        | 0  =>
            console.log "_____modelName_____ '#{mongo-query}' not found"
            reply '_____modelName_____.not-found', query
        | _  =>
            _____modelName_____ = _____modelName_____s[0]
            mongo-to-id _____modelName_____
            console.log "reading _____modelName_____ #{_____modelName_____.id}"
            reply '_____modelName_____.details', _____modelName_____


  # Updates the given _____modelName_____ object,
  # identified by its 'id' attribute
  '_____modelName_____.update': (_____modelName_____-data, {reply}) ->
    try
      id = new ObjectID _____modelName_____-data.id
    catch
      console.log "the given query (#{_____modelName_____-data}) contains an invalid id"
      return reply '_____modelName_____.not-found', id: _____modelName_____-data.id
    delete _____modelName_____-data.id
    collection.update-one {_id: id}, {$set: _____modelName_____-data}, N (result) ->
        | result.modified-count is 0  =>
            console.log "_____modelName_____ '#{id}' not updated because it doesn't exist"
            return reply '_____modelName_____.not-found'
        | _  =>
            collection.find(_id: id).to-array N (_____modelName_____s) ->
              _____modelName_____ = _____modelName_____s[0]
              mongo-to-id _____modelName_____
              console.log "updating _____modelName_____ #{_____modelName_____.id}"
              reply '_____modelName_____.updated', _____modelName_____


  '_____modelName_____.delete': (query, {reply}) ->
    try
      id = new ObjectID query.id
    catch
      console.log "the given query (#{query}) contains an invalid id"
      return reply '_____modelName_____.not-found', id: query.id
    collection.find(_id: id).to-array N (_____modelName_____s) ->
      | _____modelName_____s.length is 0  =>
          console.log "_____modelName_____ '#{id}' not deleted because it doesn't exist"
          return reply '_____modelName_____.not-found', query
      _____modelName_____ = _____modelName_____s[0]
      mongo-to-id _____modelName_____
      collection.delete-one _id: id, N (result) ->
        if result.deleted-count is 0
          console.log "_____modelName_____ '#{id}' not deleted because it doesn't exist"
          return reply '_____modelName_____.not-found', query
        console.log "deleting _____modelName_____ #{_____modelName_____.id}"
        reply '_____modelName_____.deleted', _____modelName_____


  '_____modelName_____.list': (_, {reply}) ->
    collection.find({}).to-array N (_____modelName_____s) ->
      mongo-to-ids _____modelName_____s
      console.log "listing _____modelName_____s: #{_____modelName_____s.length} found"
      reply '_____modelName_____.listing', _____modelName_____s



# Helpers

function id-to-mongo query
  result = {[k,v] for k,v of query}
  if result.id
    result._id = new ObjectID result.id
    delete result.id
  result


function mongo-to-id entry
  entry.id = entry._id
  delete entry._id
  entry


function mongo-to-ids entries
  for entry in entries
    mongo-to-id entry
