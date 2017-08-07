/**
 * Test case for ababel.
 * Runs with mocha.
 */
'use strict'

const ababel = require('../lib/ababel.js')
const assert = require('assert')
const co = require('co')

describe('ababel', function () {
  this.timeout(3000)

  before(() => co(function * () {

  }))

  after(() => co(function * () {

  }))

  it('Ababel', () => co(function * () {
    yield ababel(
      'mock-react-jsx/*.jsx',
      {
        cwd: `${__dirname}/../misc/mocks`,
        out: `${__dirname}/../tmp/testing-react-compiled`,
        presets: [ 'es2015', 'react' ],
        ext: [ '.jsx' ]
      }
    )
  }))
})

/* global describe, before, after, it */
