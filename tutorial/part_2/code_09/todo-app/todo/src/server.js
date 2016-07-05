const clone = require('clone'),
      env = require('get-env')('test'),
      {MongoClient, ObjectID} = require('mongodb'),
      N = require('nitroglycerin')


var collection = null

module.exports = {

  beforeAll: (done) => {
    const mongoDbName = `exosphere-todo-service-${env}`
    MongoClient.connect(`mongodb://localhost:27017/${mongoDbName}`, N( (mongoDb) => {
      collection = mongoDb.collection('todos')
      console.log(`MongoDB '${mongoDbName}' connected`)
      done()
    }))
  },


  // Creates a new todo object with the given data
  'todo.create': (todoData, {reply}) => {
    collection.insertOne(todoData, (err, result) => {
      if (err) {
        console.log(`Error creating todo: ${err}`)
        return reply('todo.not-created', { error: err })
      }

      console.log(`created todo '${result.ops[0]._id}'`)
      reply('todo.created', mongoToId(result.ops[0]))
    })
  },


  'todo.create-many': (todos, {reply}) => {
    collection.insert(todos, (err, result) => {
      if (err) {
        return reply('todo.not-created-many', { error: err })
      }

      console.log('creating todos: ', todos)
      reply('todo.created-many', { count: result.insertedCount })
    })
  },


  'todo.read': (query, {reply}) => {
    var mongoQuery = {}
    try {
      mongoQuery = idToMongo(query)
    } catch (e) {
      console.log(`the given query (${query}) contains an invalid id`)
      return reply('todo.not-found', query)
    }

    collection.find(mongoQuery).toArray(N( (todos) => {
      if (todos.length === 0) {
        console.log(`todo '${mongoQuery}' not found`)
        return reply('todo.not-found', query)
      }

      const todo = todos[0]
      mongoToId(todo)
      console.log(`reading todo (${todo.id})`)
      reply('todo.details', todo)
    }))
  },


  // Updates the given todo object,
  // identified by its 'id' attribute
  'todo.update': (todoData, {reply}) => {
    var id = null
    try {
      id = new ObjectID(todoData.id)
    } catch (e) {
      console.log(`the given query (${todoData}) contains an invalid id`)
      return reply('todo.not-found', { id: todoData.id })
    }
    delete todoData.id
    collection.updateOne({ _id: id }, {$set: todoData}, N( (result) => {
      if (result.modifiedCount === 0) {
        console.log(`todo '${id}' not updated because it doesn't exist`)
        return reply('todo.not-found')
      }

      collection.find({ _id: id }).toArray(N( (todos) => {
        const todo = todos[0]
        mongoToId(todo)
        console.log(`updating todo (${todo.id})`)
        reply('todo.updated', todo)
      }))
    }))
  },


  'todo.delete': (query, {reply}) => {
    var id = ''
    try {
      id = new ObjectID(query.id)
    } catch (e) {
      console.log(`the given query (${query}) contains an invalid id`)
      return reply('todo.not-found', { id: query.id })
    }

    collection.find({ _id: id }).toArray(N( (todos) => {
      if (todos.length === 0) {
        console.log(`todo '${id}' not deleted because it doesn't exist`)
        return reply('todo.not-found', query)
      }

      const todo = todos[0]
      mongoToId(todo)
      collection.deleteOne({ _id: id }, N( (result) => {
        if (result.deletedCount === 0) {
          console.log(`todo '${id}' not deleted because it doesn't exist`)
          return reply('todo.not-found', query)
        }

        console.log(`deleting todo ${todo.id}`)
        reply('todo.deleted', todo)
      }))
    }))
  },


  'todo.list': (_, {reply}) => {
    collection.find({}).toArray(N( (todos) => {
      mongoToIds(todos)
      console.log(`listing todos: ${todos.length} found`)
      reply('todo.listing', todos)
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
