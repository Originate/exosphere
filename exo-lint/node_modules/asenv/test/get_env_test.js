/**
 * Test case for getEnv.
 * Runs with mocha.
 */
'use strict'

const getEnv = require('../lib/get_env.js')
const assert = require('assert')
const co = require('co')

describe('get-env', function () {
  this.timeout(3000)

  before(() => co(function * () {
  }))

  after(() => co(function * () {

  }))

  it('Get env', () => co(function * () {
    process.env.NODE_ENV = 'fuge'
    assert.equal(getEnv(), 'fuge')
  }))
})

/* global describe, before, after, it */
