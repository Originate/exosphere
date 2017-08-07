/**
 * Test case for create.
 * Runs with mocha.
 */
'use strict'

const create = require('../lib/create.js')
const assert = require('assert')
const co = require('co')

describe('create', function () {
  this.timeout(3000)

  before(() => co(function * () {

  }))

  after(() => co(function * () {

  }))

  it('Create', () => co(function * () {
    let filename = `${__dirname}/../tmp/foo/bar.json`
    let store = create(filename)
    assert.ok(store)
    assert.equal(store.storage.filename, filename)
  }))
})

/* global describe, before, after, it */
