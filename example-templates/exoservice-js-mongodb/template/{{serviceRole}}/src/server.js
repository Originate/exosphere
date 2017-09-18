const clone = require('clone'),
      env = require('get-env')('test'),
      {MongoClient, ObjectID} = require('mongodb'),
      N = require('nitroglycerin'),
      {bootstrap} = require('exoservice')


var collection = null

bootstrap({

  beforeAll: (done) => {
    const mongoDbName = `exosphere-{{serviceRole}}-${env}`
    MongoClient.connect(`mongodb://${process.env.MONGO}:27017/${mongoDbName}`, N( (mongoDb) => {
      collection = mongoDb.collection('{{modelName}}s')
      console.log(`MongoDB '${mongoDbName}' connected`)
      done()
    }))
  },


  // Creates a new {{modelName}} object with the given data
  '{{modelName}}.create': ({{modelName}}Data, {reply}) => {
    collection.insertOne({{modelName}}Data, (err, result) => {
      if (err) {
        console.log(`Error creating {{modelName}}: ${err}`)
        return reply('{{modelName}}.not-created', { error: err })
      }

      console.log(`created {{modelName}} '${result.ops[0]._id}'`)
      reply('{{modelName}}.created', mongoToId(result.ops[0]))
    })
  },


  '{{modelName}}.create-many': ({{modelName}}s, {reply}) => {
    collection.insert({{modelName}}s, (err, result) => {
      if (err) {
        return reply('{{modelName}}.not-created-many', { error: err })
      }

      console.log('creating {{modelName}}s: ', {{modelName}}s)
      reply('{{modelName}}.created-many', { count: result.insertedCount })
    })
  },


  '{{modelName}}.read': (query, {reply}) => {
    var mongoQuery = {}
    try {
      mongoQuery = idToMongo(query)
    } catch (e) {
      console.log(`the given query (${query}) contains an invalid id`)
      return reply('{{modelName}}.not-found', query)
    }

    collection.find(mongoQuery).toArray(N( ({{modelName}}s) => {
      if ({{modelName}}s.length === 0) {
        console.log(`{{modelName}} '${mongoQuery}' not found`)
        return reply('{{modelName}}.not-found', query)
      }

      const {{modelName}} = {{modelName}}s[0]
      mongoToId({{modelName}})
      console.log(`reading {{modelName}} (${ {{modelName}}.id })`)
      reply('{{modelName}}.details', {{modelName}})
    }))
  },


  // Updates the given {{modelName}} object,
  // identified by its 'id' attribute
  '{{modelName}}.update': ({{modelName}}Data, {reply}) => {
    var id = null
    try {
      id = new ObjectID({{modelName}}Data.id)
    } catch (e) {
      console.log(`the given query (${ {{modelName}}Data }) contains an invalid id`)
      return reply('{{modelName}}.not-found', { id: {{modelName}}Data.id })
    }
    delete {{modelName}}Data.id
    collection.updateOne({ _id: id }, {$set: {{modelName}}Data}, N( (result) => {
      if (result.modifiedCount === 0) {
        console.log(`{{modelName}} '${id}' not updated because it doesn't exist`)
        return reply('{{modelName}}.not-found')
      }

      collection.find({ _id: id }).toArray(N( ({{modelName}}s) => {
        const {{modelName}} = {{modelName}}s[0]
        mongoToId({{modelName}})
        console.log(`updating {{modelName}} (${ {{modelName}}.id })`)
        reply('{{modelName}}.updated', {{modelName}})
      }))
    }))
  },


  '{{modelName}}.delete': (query, {reply}) => {
    var id = ''
    try {
      id = new ObjectID(query.id)
    } catch (e) {
      console.log(`the given query (${query}) contains an invalid id`)
      return reply('{{modelName}}.not-found', { id: query.id })
    }

    collection.find({ _id: id }).toArray(N( ({{modelName}}s) => {
      if ({{modelName}}s.length === 0) {
        console.log(`{{modelName}} '${id}' not deleted because it doesn't exist`)
        return reply('{{modelName}}.not-found', query)
      }

      const {{modelName}} = {{modelName}}s[0]
      mongoToId({{modelName}})
      collection.deleteOne({ _id: id }, N( (result) => {
        if (result.deletedCount === 0) {
          console.log(`{{modelName}} '${id}' not deleted because it doesn't exist`)
          return reply('{{modelName}}.not-found', query)
        }

        console.log(`deleting {{modelName}} ${ {{modelName}}.id }`)
        reply('{{modelName}}.deleted', {{modelName}})
      }))
    }))
  },


  '{{modelName}}.list': (_, {reply}) => {
    collection.find({}).toArray(N( ({{modelName}}s) => {
      mongoToIds({{modelName}}s)
      console.log(`listing {{modelName}}s: ${ {{modelName}}s.length } found`)
      reply('{{modelName}}.listing', {{modelName}}s)
    }))
  }

})


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
