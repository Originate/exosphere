/**
 * Test case for convertString.
 * Runs with mocha.
 */
'use strict'

const convertString = require('../lib/converting/convert_string.js')
const co = require('co')
const assert = require('assert')
before(() => co(function * () {

}))

after(() => co(function * () {

}))

it('Convert string', () => co(function * () {
  let converted = convertString('foo bar baz bar', {
    bar: 'quz'
  })
  assert.equal(converted, 'foo quz baz quz')
}))

/* global describe, before, after, it */
