/**
 * Test case for filedel.
 * Runs with mocha.
 */
'use strict'

const filedel = require('../shim/node/filedel.js')
const path = require('path')
const fs = require('fs')
const co = require('co')
const mkdirp = require('mkdirp')
const assert = require('assert')

let tmpDir = path.resolve(__dirname, '../tmp')

describe('filedel', () => {
  before(() => {
    mkdirp.sync(tmpDir)
  })

  after(() => {

  })

  it('Unlink a file.', () => co(function * () {
    let filename = path.resolve(tmpDir, 'work_file_to_unlink.txt')
    fs.writeFileSync(filename, 'foo')
    assert.ok(fs.existsSync(filename))
    yield filedel(filename, {
      force: true
    })
    assert.ok(!fs.existsSync(filename))
    yield filedel(filename)
    assert.ok(!fs.existsSync(filename))
  }))

  it('Unlink dir', () => co(function * () {
    let dirname = `${tmpDir}/hoge/un-linking-dir`
    mkdirp.sync(dirname)
    yield filedel.recursive(dirname)
  }))

  it('Try to delete dir.', () => co(function * () {
    let dirname = path.resolve(tmpDir, 'work_dir_to_unlink')
    mkdirp.sync(dirname)
    try {
      yield filedel(dirname)
    } catch (err) {
      assert.ok(!!err)
    }
  }))
})

/* global describe, before, after, it */
