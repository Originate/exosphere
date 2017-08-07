/**
 * Test case for convertFile.
 * Runs with mocha.
 */
'use strict'

const convertFile = require('../lib/converting/convert_file.js')
const mkdirp = require('mkdirp')
const co = require('co')
const assert = require('assert')

const tmpDir = __dirname + '/../tmp';

before(() => co(function * () {
  mkdirp.sync(tmpDir)
}))

after(() => co(function * () {

}))

it('Convert file', () => co(function * () {
  let src = String(__filename)
  let dest = tmpDir + '/foo/bar/testing-converted.txt';
  yield convertFile(src, dest, {})
}))

/* global describe, before, after, it */
