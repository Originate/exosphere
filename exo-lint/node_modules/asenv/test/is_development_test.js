/**
 * Test case for isDevelopment.
 * Runs with mocha.
 */
'use strict'

const isDevelopment = require('../lib/is_development.js')
const assert = require('assert')
const co = require('co')

describe('is-development', function () {
  this.timeout(3000)

  before(() => co(function * () {

  }))

  after(() => co(function * () {

  }))

  it('Is development', () => co(function * () {
    process.env.NODE_ENV = 'development'
    assert.ok(isDevelopment())
    process.env.NODE_ENV = 'production'
    assert.ok(!isDevelopment())
  }))
})

/* global describe, before, after, it */
