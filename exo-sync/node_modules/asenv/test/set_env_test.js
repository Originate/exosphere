/**
 * Test case for setEnv.
 * Runs with mocha.
 */
'use strict'

const setEnv = require('../lib/set_env.js')
const assert = require('assert')
const co = require('co')

describe('set-env', function () {
  this.timeout(3000)

  before(() => co(function * () {

  }))

  after(() => co(function * () {

  }))

  it('Set env', () => co(function * () {
    setEnv('hoge')
    assert.equal(process.env.NODE_ENV, 'hoge')
  }))
})

/* global describe, before, after, it */
