/**
 * Test case for isProduction.
 * Runs with mocha.
 */
'use strict'

const isProduction = require('../lib/is_production.js')
const assert = require('assert')
const co = require('co')

describe('is-production', function () {
  this.timeout(3000)

  before(() => co(function * () {

  }))

  after(() => co(function * () {

  }))

  it('Is production', () => co(function * () {
    process.env.NODE_ENV = 'production'
    assert.ok(isProduction())
    process.env.NODE_ENV = 'test'
    assert.ok(!isProduction())
  }))
})

/* global describe, before, after, it */
