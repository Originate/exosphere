/**
 * Test case for tmplify.
 * Runs with mocha.
 */
'use strict'

const tmplify = require('../lib/tmplify.js')
const mkdirp = require('mkdirp')

const tmpDir = __dirname + '/../tmp';
const co = require('co')
const assert = require('assert')
const fs = require('fs')

describe('tmplify', function () {
  before(() => co(function * () {
    mkdirp.sync(tmpDir)
  }))

  it('Tmplify', () => co(function * () {
    let srcDir = __dirname + '/../doc/mocks/mock-app'
    let destDir = tmpDir + '/testing-tmpl/mock-app-tmpl'
    yield tmplify(srcDir, destDir, {
      data: {
        "name": "my-awesome-app",
        "description": "This is an example for the app templates."
      }
    })
    assert.ok(
      fs.existsSync(destDir)
    )
    let data = JSON.parse(fs.readFileSync(
      `${tmpDir}/testing-tmpl/mock-app-tmpl/package.json.tmpl`
    ).toString())
    assert.deepEqual(data, {
        name: '_____name_____',
        version: '1.0.0',
        description: '_____description_____'
      }
    )
  }))

})

/* global describe, before, after, it */
