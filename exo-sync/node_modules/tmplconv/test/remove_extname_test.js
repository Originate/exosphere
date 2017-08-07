/**
 * Test case for removeExtname.
 * Runs with mocha.
 */
'use strict'

const removeExtname = require('../lib/naming/remove_extname.js')
const co = require('co')
const assert = require('assert')

it('Remove extname', () => co(function * () {
  assert.equal(removeExtname('foo/bar/baz.txt.tmpl', '.tmpl'), 'foo/bar/baz.txt')
  assert.equal(removeExtname('foo/bar/baz.txt', '.tmpl'), 'foo/bar/baz.txt')
}))

/* global describe, before, after, it */
