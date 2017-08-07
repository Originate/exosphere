/**
 * Test case for akvStatus.
 * Runs with mocha.
 */
'use strict'

const AkvStatus = require('../lib/akv_status.js')
const assert = require('assert')
const co = require('co')
const mkdirp = require('mkdirp')
const fs = require('fs')
const path = require('path')

describe('akv-status', function () {
  this.timeout(3000)

  before(() => co(function * () {

  }))

  after(() => co(function * () {

  }))

  it('Akv status', () => co(function * () {
    let tmplDir = `${__dirname}/../tmp/testing-status-file`
    mkdirp.sync(tmplDir)
    let [filename01, filename02] = [
      path.resolve(tmplDir, 'hoge01.txt'),
      path.resolve(tmplDir, 'hoge02.txt')
    ]
    fs.writeFileSync(filename01, '01-' + new Date())
    fs.writeFileSync(filename02, '02-' + new Date())
    let store = new AkvStatus(`${tmplDir}/testing.status.json`)
    yield store.saveStatus([ filename01, filename02 ])
    fs.writeFileSync(filename01, '01-updated-' + new Date())
    {
      let filtered = yield store.filterStatusUnknown([ filename01, filename02 ])
      assert.deepEqual(filtered, [ filename01 ])
    }
    fs.writeFileSync(filename02, '02-updated-' + new Date())
    {
      let filtered = yield store.filterStatusUnknown([ filename01, filename02 ])
      assert.deepEqual(filtered, [ filename01, filename02 ])
    }
  }))
})

/* global describe, before, after, it */
