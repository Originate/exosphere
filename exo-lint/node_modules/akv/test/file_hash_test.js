/**
 * Test case for fileHash.
 * Runs with mocha.
 */
'use strict'

const fileHash = require('../lib/file/file_hash.js')
const assert = require('assert')
const co = require('co')

describe('file-hash', function () {
  this.timeout(3000)

  before(() => co(function * () {

  }))

  after(() => co(function * () {

  }))

  it('File hash', () => co(function * () {
    {
      let hash = yield fileHash(__filename)
      assert.ok(hash)
    }
  }))
})

/* global describe, before, after, it */
