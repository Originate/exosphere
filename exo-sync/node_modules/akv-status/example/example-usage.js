'use strict'

const akvStatus = require('akv-status')
const co = require('co')

co(function * () {
  let store = akvStatus('tmp/files.status.json')
  let patterns = [ 'src/**/*.js' ]
  yield store.saveStatus(patterns)

  process.on('message', () => co(function * () {
    let changed = yield store.filterStatusUnknown(patterns)
    if (changed.length > 0) {
      /* ... */
    }
  }))
}).catch((err) => console.error(err))
