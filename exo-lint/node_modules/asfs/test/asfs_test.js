/**
 * Test case for asfs.
 * Runs with mocha.
 */
'use strict'

const asfs = require('../lib/asfs.js')
const assert = require('assert')
const co = require('co')

describe('asfs', function () {
  this.timeout(3000)

  before(() => co(function * () {

  }))

  after(() => co(function * () {

  }))

  it('Asfs', () => co(function * () {
    yield asfs.mkdirpAsync(`${__dirname}/../tmp/foo`)
    yield asfs.writeFileAsync(`${__dirname}/../tmp/foo/bar.txt`, 'This is bar')
    let exists = yield asfs.existsAsync(`${__dirname}/../tmp/foo/bar.txt`)
    assert.ok(exists)
    let content = yield asfs.readFileAsync(`${__dirname}/../tmp/foo/bar.txt`)
    assert.equal(content.toString(), 'This is bar')

    // With encode
    {
      let content = yield asfs.readFileAsync(`${__dirname}/../tmp/foo/bar.txt`, 'base64')
      assert.equal(content.toString(), 'VGhpcyBpcyBiYXI=')
    }

    {
      let stat = yield asfs.statAsync(`${__dirname}/../tmp/foo/bar.txt`, 'base64')
      console.log(stat)
      assert.equal(stat.size, 11)
    }
    {
      let filenames = yield asfs.readdirAsync(__dirname)
      assert.ok(filenames)
    }
  }))
})

/* global describe, before, after, it */
