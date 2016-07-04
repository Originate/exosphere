const clone = require('clone'),
      env = require('get-env')('test'),
      {MongoClient, ObjectID} = require('mongodb'),
      N = require('nitroglycerin')


var collection = null

module.exports = {

  beforeAll: (done) => {
    const mongoDbName = `exosphere-_____serviceName_____-service-${env}`
    MongoClient.connect(`mongodb://localhost:27017/${mongoDbName}`, N( (mongoDb) => {
      collection = mongoDb.collection('_____serviceName_____s')
      console.log(`MongoDB '${mongoDbName}' connected`)
      done()
    }))
  },


  // Creates a new _____serviceName_____ object with the given data
  '_____serviceName_____.create': (_____serviceName_____Data, {reply}) => {
    collection.insertOne(_____serviceName_____Data, (err, result) => {
      if (err) {
        console.log(`Error creating _____serviceName_____: ${err}`)
        return reply('_____serviceName_____.not-created', { error: err })
      }

      console.log(`created _____serviceName_____ '${result.ops[0]._id}'`)
      reply('_____serviceName_____.created', mongoToId(result.ops[0]))
    })
  },


  '_____serviceName_____.create-many': (_____serviceName_____s, {reply}) => {
    collection.insert(_____serviceName_____s, (err, result) => {
      if (err) {
        return reply('_____serviceName_____.not-created-many', { error: err })
      }

      console.log('creating _____serviceName_____s: ', _____serviceName_____s)
      reply('_____serviceName_____.created-many', { count: result.insertedCount })
    })
  },


  '_____serviceName_____.read': (query, {reply}) => {
    var mongoQuery = {}
    try {
      mongoQuery = idToMongo(query)
    } catch (e) {
      console.log(`the given query (${query}) contains an invalid id`)
      return reply('_____serviceName_____.not-found', query)
    }

    collection.find(mongoQuery).toArray(N( (_____serviceName_____s) => {
      if (_____serviceName_____s.length === 0) {
        console.log(`_____serviceName_____ '${mongoQuery}' not found`)
        return reply('_____serviceName_____.not-found', query)
      }

      const _____serviceName_____ = _____serviceName_____s[0]
      mongoToId(_____serviceName_____)
      console.log(`reading _____serviceName_____ (${_____serviceName_____.id})`)
      reply('_____serviceName_____.details', _____serviceName_____)
    }))
  },


  // Updates the given _____serviceName_____ object,
  // identified by its 'id' attribute
  '_____serviceName_____.update': (_____serviceName_____Data, {reply}) => {
    var id = null
    try {
      id = new ObjectID(_____serviceName_____Data.id)
    } catch (e) {
      console.log(`the given query (${_____serviceName_____Data}) contains an invalid id`)
      return reply('_____serviceName_____.not-found', { id: _____serviceName_____Data.id })
    }
    delete _____serviceName_____Data.id
    collection.updateOne({ _id: id }, {$set: _____serviceName_____Data}, N( (result) => {
      if (result.modifiedCount === 0) {
        console.log(`_____serviceName_____ '${id}' not updated because it doesn't exist`)
        return reply('_____serviceName_____.not-found')
      }

      collection.find({ _id: id }).toArray(N( (_____serviceName_____s) => {
        const _____serviceName_____ = _____serviceName_____s[0]
        mongoToId(_____serviceName_____)
        console.log(`updating _____serviceName_____ (${_____serviceName_____.id})`)
        reply('_____serviceName_____.updated', _____serviceName_____)
      }))
    }))
  },


  '_____serviceName_____.delete': (query, {reply}) => {
    var id = ''
    try {
      id = new ObjectID(query.id)
    } catch (e) {
      console.log(`the given query (${query}) contains an invalid id`)
      return reply('_____serviceName_____.not-found', { id: query.id })
    }

    collection.find({ _id: id }).toArray(N( (_____serviceName_____s) => {
      if (_____serviceName_____s.length === 0) {
        console.log(`_____serviceName_____ '${id}' not deleted because it doesn't exist`)
        return reply('_____serviceName_____.not-found', query)
      }

      const _____serviceName_____ = _____serviceName_____s[0]
      mongoToId(_____serviceName_____)
      collection.deleteOne({ _id: id }, N( (result) => {
        if (result.deletedCount === 0) {
          console.log(`_____serviceName_____ '${id}' not deleted because it doesn't exist`)
          return reply('_____serviceName_____.not-found', query)
        }

        console.log(`deleting _____serviceName_____ ${_____serviceName_____.id}`)
        reply('_____serviceName_____.deleted', _____serviceName_____)
      }))
    }))
  },


  '_____serviceName_____.list': (_, {reply}) => {
    collection.find({}).toArray(N( (_____serviceName_____s) => {
      mongoToIds(_____serviceName_____s)
      console.log(`listing _____serviceName_____s: ${_____serviceName_____s.length} found`)
      reply('_____serviceName_____.listing', _____serviceName_____s)
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
