/**
 * Test case for storage.
 * Runs with mocha.
 */
'use strict'

const Storage = require('../lib/storage.js')
const assert = require('assert')
const co = require('co')
const fs = require('fs')

describe('storage', function () {
  this.timeout(3000)

  before(() => co(function * () {

  }))

  after(() => co(function * () {

  }))

  it('Storage', () => co(function * () {
    let filename = `${__dirname}/../tmp/testing-storage/storage01.json`
    let storage = new Storage(
      filename,
      { interval: 100 }
    )
    {
      yield storage.write({ foo: 'bar' })
      let data = yield storage.read()
      assert.deepEqual(data, { foo: 'bar' })
    }
    {
      yield storage.write({ foo: 'baz' })
      let data = yield storage.read()
      assert.deepEqual(data, { foo: 'baz' })
    }
    yield storage.flush()
    assert.ok(fs.existsSync(filename))
    yield storage.purge()
    yield storage.purge()
    yield storage.purge()
    assert.ok(!fs.existsSync(filename))
  }))
})

/* global describe, before, after, it */
