const clone = require('clone'),
      env = require('get-env')('test'),
      {MongoClient, ObjectID} = require('mongodb'),
      N = require('nitroglycerin')


var collection = null

module.exports = {

  beforeAll: (done) => {
    const mongoDbName = `exosphere-_____serviceName_____-${env}`
    MongoClient.connect(`mongodb://localhost:27017/${mongoDbName}`, N( (mongoDb) => {
      collection = mongoDb.collection('_____modelName_____s')
      console.log(`MongoDB '${mongoDbName}' connected`)
      done()
    }))
  },


  // Creates a new _____modelName_____ object with the given data
  '_____modelName_____.create': (_____modelName_____Data, {reply}) => {
    collection.insertOne(_____modelName_____Data, (err, result) => {
      if (err) {
        console.log(`Error creating _____modelName_____: ${err}`)
        return reply('_____modelName_____.not-created', { error: err })
      }

      console.log(`created _____modelName_____ '${result.ops[0]._id}'`)
      reply('_____modelName_____.created', mongoToId(result.ops[0]))
    })
  },


  '_____modelName_____.create-many': (_____modelName_____s, {reply}) => {
    collection.insert(_____modelName_____s, (err, result) => {
      if (err) {
        return reply('_____modelName_____.not-created-many', { error: err })
      }

      console.log('creating _____modelName_____s: ', _____modelName_____s)
      reply('_____modelName_____.created-many', { count: result.insertedCount })
    })
  },


  '_____modelName_____.read': (query, {reply}) => {
    var mongoQuery = {}
    try {
      mongoQuery = idToMongo(query)
    } catch (e) {
      console.log(`the given query (${query}) contains an invalid id`)
      return reply('_____modelName_____.not-found', query)
    }

    collection.find(mongoQuery).toArray(N( (_____modelName_____s) => {
      if (_____modelName_____s.length === 0) {
        console.log(`_____modelName_____ '${mongoQuery}' not found`)
        return reply('_____modelName_____.not-found', query)
      }

      const _____modelName_____ = _____modelName_____s[0]
      mongoToId(_____modelName_____)
      console.log(`reading _____modelName_____ (${_____modelName_____.id})`)
      reply('_____modelName_____.details', _____modelName_____)
    }))
  },


  // Updates the given _____modelName_____ object,
  // identified by its 'id' attribute
  '_____modelName_____.update': (_____modelName_____Data, {reply}) => {
    var id = null
    try {
      id = new ObjectID(_____modelName_____Data.id)
    } catch (e) {
      console.log(`the given query (${_____modelName_____Data}) contains an invalid id`)
      return reply('_____modelName_____.not-found', { id: _____modelName_____Data.id })
    }
    delete _____modelName_____Data.id
    collection.updateOne({ _id: id }, {$set: _____modelName_____Data}, N( (result) => {
      if (result.modifiedCount === 0) {
        console.log(`_____modelName_____ '${id}' not updated because it doesn't exist`)
        return reply('_____modelName_____.not-found')
      }

      collection.find({ _id: id }).toArray(N( (_____modelName_____s) => {
        const _____modelName_____ = _____modelName_____s[0]
        mongoToId(_____modelName_____)
        console.log(`updating _____modelName_____ (${_____modelName_____.id})`)
        reply('_____modelName_____.updated', _____modelName_____)
      }))
    }))
  },


  '_____modelName_____.delete': (query, {reply}) => {
    var id = ''
    try {
      id = new ObjectID(query.id)
    } catch (e) {
      console.log(`the given query (${query}) contains an invalid id`)
      return reply('_____modelName_____.not-found', { id: query.id })
    }

    collection.find({ _id: id }).toArray(N( (_____modelName_____s) => {
      if (_____modelName_____s.length === 0) {
        console.log(`_____modelName_____ '${id}' not deleted because it doesn't exist`)
        return reply('_____modelName_____.not-found', query)
      }

      const _____modelName_____ = _____modelName_____s[0]
      mongoToId(_____modelName_____)
      collection.deleteOne({ _id: id }, N( (result) => {
        if (result.deletedCount === 0) {
          console.log(`_____modelName_____ '${id}' not deleted because it doesn't exist`)
          return reply('_____modelName_____.not-found', query)
        }

        console.log(`deleting _____modelName_____ ${_____modelName_____.id}`)
        reply('_____modelName_____.deleted', _____modelName_____)
      }))
    }))
  },


  '_____modelName_____.list': (_, {reply}) => {
    collection.find({}).toArray(N( (_____modelName_____s) => {
      mongoToIds(_____modelName_____s)
      console.log(`listing _____modelName_____s: ${_____modelName_____s.length} found`)
      reply('_____modelName_____.listing', _____modelName_____s)
    }))
  }

}


// Helpers

function idToMongo(query) {
  const result = clone(query)
  if (result.id) {
    result._id = new ObjectID(result.id)
    delete result.id
  }
  return result
}


function mongoToId(entry) {
  entry.id = entry._id
  delete entry._id
  return entry
}


function mongoToIds(entries) {
  for (entry in entries) {
    mongoToId(entry)
  }
}
