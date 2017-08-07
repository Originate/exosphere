/**
 * Test for index.js
 * Runs with nodeunit.
 */

'use strict'

const index = require('../lib/index')
const assert = require('assert')

describe('index', () => {
  it('Eval properties.', () => {
    assert.ok(index)
    Object.keys(index).forEach((key) => {
      assert.ok(key)
      assert.ok(index[ key ])
    })
  })
})

/* global describe, before, after, it */