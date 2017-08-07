'use strict'

const akv = require('akv')
const co = require('co')

co(function * () {
  let storage = akv('tmp/my-storage.json')
  // Set key value
  yield storage.set('foo', 'bar')

  // Get key value
  let foo = yield storage.get('foo')
  console.log(foo) // => bar
  // Delete by key
  yield storage.del('foo')
}).catch((err) => console.error(err))
