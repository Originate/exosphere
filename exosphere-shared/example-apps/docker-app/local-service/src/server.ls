require! {
  'mongodb' : {MongoClient, ObjectID}
  'nitroglycerin' : N
}
env = require('get-env')('test')


collection = null

module.exports =

  before-all: (done) ->
    mongo-db-name = "exosphere-local-service-#{env}"
    mongo-address = if env is \test then \localhost else process.env.MONGO
    MongoClient.connect "mongodb://#{mongo-address}:27017/#{mongo-db-name}", N (mongo-db) ->
      collection := mongo-db.collection 'tests'
      console.log "MongoDB '#{mongo-db-name}' connected"
      done!


  # Creates a new test object with the given data
  'test.create': (test-data, {reply}) ->
    collection.insert-one test-data, (err, result) ->
      | err  =>
          console.log "Error creating test: #{err}"
          reply 'test.not-created', error: err
      | _  =>
          console.log "creating test"
          reply 'test.created', mongo-to-id(result.ops[0])


  'test.create-many': (tests, {reply}) ->
    collection.insert tests, (err, result) ->
      | err  =>  return reply 'test.not-created-many', error: err
      reply 'test.created-many', count: result.inserted-count


  'test.read': (query, {reply}) ->
    try
      mongo-query = id-to-mongo query
    catch
      console.log "the given query (#{query}) contains an invalid id"
      return reply 'test.not-found', query
    collection.find(mongo-query).to-array N (tests) ->
      switch tests.length
        | 0  =>
            console.log "test '#{mongo-query}' not found"
            reply 'test.not-found', query
        | _  =>
            test = tests[0]
            mongo-to-id test
            console.log "reading test #{test.id}"
            reply 'test.details', test


  # Updates the given test object,
  # identified by its 'id' attribute
  'test.update': (test-data, {reply}) ->
    try
      id = new ObjectID test-data.id
    catch
      console.log "the given query (#{test-data}) contains an invalid id"
      return reply 'test.not-found', id: test-data.id
    delete test-data.id
    collection.update-one {_id: id}, {$set: test-data}, N (result) ->
        | result.modified-count is 0  =>
            console.log "test '#{id}' not updated because it doesn't exist"
            return reply 'test.not-found'
        | _  =>
            collection.find(_id: id).to-array N (tests) ->
              test = tests[0]
              mongo-to-id test
              console.log "updating test #{test.id}"
              reply 'test.updated', test


  'test.delete': (query, {reply}) ->
    try
      id = new ObjectID query.id
    catch
      console.log "the given query (#{query}) contains an invalid id"
      return reply 'test.not-found', id: query.id
    collection.find(_id: id).to-array N (tests) ->
      | tests.length is 0  =>
          console.log "test '#{id}' not deleted because it doesn't exist"
          return reply 'test.not-found', query
      test = tests[0]
      mongo-to-id test
      collection.delete-one _id: id, N (result) ->
        if result.deleted-count is 0
          console.log "test '#{id}' not deleted because it doesn't exist"
          return reply 'test.not-found', query
        console.log "deleting test #{test.id}"
        reply 'test.deleted', test


  'test.list': (_, {reply}) ->
    collection.find({}).to-array N (tests) ->
      mongo-to-ids tests
      console.log "listing tests: #{tests.length} found"
      reply 'test.listing', tests



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
