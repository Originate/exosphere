/**
 * Test case for isTest.
 * Runs with mocha.
 */
'use strict'

const isTest = require('../lib/is_test.js')
const assert = require('assert')
const co = require('co')

describe('is-test', function () {
  this.timeout(3000)

  before(() => co(function * () {

  }))

  after(() => co(function * () {

  }))

  it('Is test', () => co(function * () {
    process.env.NODE_ENV = 'test'
    assert.ok(isTest())
  }))
})

/* global describe, before, after, it */
