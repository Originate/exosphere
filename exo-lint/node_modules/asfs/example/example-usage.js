'use strict'

const asfs = require('asfs')
const co = require('co')

co(function * () {
  let filename = 'foo.txt'
  let exists = yield asfs.existsAsync(filename)
  if (exists) {
    let content = yield asfs.readFileAsync(filename)
    console.log(content)
  }
})
