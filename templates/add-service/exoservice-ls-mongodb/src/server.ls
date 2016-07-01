require! {
  'mongodb' : {MongoClient, ObjectID}
  'nitroglycerin' : N
}
env = require('get-env')('test')


collection = null

module.exports =

  before-all: (done) ->
    mongo-db-name = "exosphere-_____serviceName_____-service-#{env}"
    MongoClient.connect "mongodb://localhost:27017/#{mongo-db-name}", N (mongo-db) ->
      collection := mongo-db.collection '_____serviceName_____s'
      console.log "MongoDB '#{mongo-db-name}' connected"
      done!


  # Creates a new _____serviceName_____ object with the given data
  '_____serviceName_____.create': (_____serviceName_____-data, {reply}) ->
    collection.insert-one _____serviceName_____-data, (err, result) ->
      | err  =>
          console.log "Error creating _____serviceName_____: #{err}"
          reply '_____serviceName_____.not-created', error: err
      | _  =>
          console.log "creating _____serviceName_____"
          reply '_____serviceName_____.created', mongo-to-id(result.ops[0])


  '_____serviceName_____.create-many': (_____serviceName_____s, {reply}) ->
    collection.insert _____serviceName_____s, (err, result) ->
      | err  =>  return reply '_____serviceName_____.not-created-many', error: err
      reply '_____serviceName_____.created-many', count: result.inserted-count


  '_____serviceName_____.read': (query, {reply}) ->
    try
      mongo-query = id-to-mongo query
    catch
      console.log "the given query (#{query}) contains an invalid id"
      return reply '_____serviceName_____.not-found', query
    collection.find(mongo-query).to-array N (_____serviceName_____s) ->
      switch _____serviceName_____s.length
        | 0  =>
            console.log "_____serviceName_____ '#{mongo-query}' not found"
            reply '_____serviceName_____.not-found', query
        | _  =>
            _____serviceName_____ = _____serviceName_____s[0]
            mongo-to-id _____serviceName_____
            console.log "reading _____serviceName_____ #{_____serviceName_____.id}"
            reply '_____serviceName_____.details', _____serviceName_____


  # Updates the given _____serviceName_____ object,
  # identified by its 'id' attribute
  '_____serviceName_____.update': (_____serviceName_____-data, {reply}) ->
    try
      id = new ObjectID _____serviceName_____-data.id
    catch
      console.log "the given query (#{_____serviceName_____-data}) contains an invalid id"
      return reply '_____serviceName_____.not-found', id: _____serviceName_____-data.id
    delete _____serviceName_____-data.id
    collection.update-one {_id: id}, {$set: _____serviceName_____-data}, N (result) ->
        | result.modified-count is 0  =>
            console.log "_____serviceName_____ '#{id}' not updated because it doesn't exist"
            return reply '_____serviceName_____.not-found'
        | _  =>
            collection.find(_id: id).to-array N (_____serviceName_____s) ->
              _____serviceName_____ = _____serviceName_____s[0]
              mongo-to-id _____serviceName_____
              console.log "updating _____serviceName_____ #{_____serviceName_____.id}"
              reply '_____serviceName_____.updated', _____serviceName_____


  '_____serviceName_____.delete': (query, {reply}) ->
    try
      id = new ObjectID query.id
    catch
      console.log "the given query (#{query}) contains an invalid id"
      return reply '_____serviceName_____.not-found', id: query.id
    collection.find(_id: id).to-array N (_____serviceName_____s) ->
      | _____serviceName_____s.length is 0  =>
          console.log "_____serviceName_____ '#{id}' not deleted because it doesn't exist"
          return reply '_____serviceName_____.not-found', query
      _____serviceName_____ = _____serviceName_____s[0]
      mongo-to-id _____serviceName_____
      collection.delete-one _id: id, N (result) ->
        if result.deleted-count is 0
          console.log "_____serviceName_____ '#{id}' not deleted because it doesn't exist"
          return reply '_____serviceName_____.not-found', query
        console.log "deleting _____serviceName_____ #{_____serviceName_____.id}"
        reply '_____serviceName_____.deleted', _____serviceName_____


  '_____serviceName_____.list': (_, {reply}) ->
    collection.find({}).to-array N (_____serviceName_____s) ->
      mongo-to-ids _____serviceName_____s
      console.log "listing _____serviceName_____s: #{_____serviceName_____s.length} found"
      reply '_____serviceName_____.listed', {count: _____serviceName_____s.length, _____serviceName_____s}



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
