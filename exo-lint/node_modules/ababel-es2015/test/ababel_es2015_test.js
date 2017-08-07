/**
 * Test case for ababelEs2015.
 * Runs with mocha.
 */
'use strict'

const ababelEs2015 = require('../lib/ababel_es2015.js')
const assert = require('assert')
const co = require('co')

describe('ababel-es2015', function () {
  this.timeout(3000)

  before(() => co(function * () {

  }))

  after(() => co(function * () {

  }))

  it('Ababel es2015', () => co(function * () {
    yield ababelEs2015(
      'mocks/*.js',
      {
        cwd: `${__dirname}/../misc`,
        out: `${__dirname}/../tmp/testing-compiled`,
        minified: true
      }
    )
  }))
})

/* global describe, before, after, it */
