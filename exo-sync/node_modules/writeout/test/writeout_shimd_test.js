/**
 * Test for writeout.js
 * Runs with nodeunit.
 */

'use strict'

const writeout = require('../shim/node')
const co = require('co')
const path = require('path')
const assert = require('assert')
const mkdirp = require('mkdirp')

let tmpDir = path.resolve(__dirname, '../tmp')

describe('writeout', () => {
  before(() => {
    mkdirp.sync(tmpDir)
  })

  it('Do writeout.', () => co(function * () {
    let filename = path.resolve(tmpDir, 'foo/bar/baz.txt')
    let result = yield writeout(filename, 'Oh!', {
      mkdirp: true,
      skipIfIdentical: false
    })
    assert.ok(result)
    let result2 = yield writeout(filename, 'Oh!', {
      mkdirp: false,
      skipIfIdentical: true
    })
    assert.equal(result2.skipped, true)
  }))

  it('Do writeout to force.', () => co(function * () {
    let filename = path.resolve(tmpDir, 'foo/bar/baz-readonly.txt')
    yield writeout(filename, 'This is readonly.', {
      mkdirp: true,
      skipIfIdentical: false,
      mode: '444'
    })
    yield writeout(filename, 'Override readonly file.', {
      mkdirp: false,
      skipIfIdentical: true,
      mode: '644',
      force: true
    })
  }))
})

/* global describe, before, it */
