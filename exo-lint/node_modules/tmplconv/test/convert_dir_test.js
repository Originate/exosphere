/**
 * Test case for convertDir.
 * Runs with mocha.
 */
'use strict'

const convertDir = require('../lib/converting/convert_dir.js')
const mkdirp = require('mkdirp')
const co = require('co')
const assert = require('assert')

const tmpDir = __dirname + '/../tmp';

before(() => co(function * () {
  mkdirp.sync(tmpDir)
}))

it('Convert dir', () => co(function * () {
  yield  convertDir(__dirname, tmpDir + '/baz', {
    pattern: '*.*'
  })
}))

/* global describe, before, after, it */
