/**
 * Test case for akv.
 * Runs with mocha.
 */
'use strict'

const AKV = require('../lib/akv.js')
const assert = require('assert')
const asleep = require('asleep')
const fs = require('fs')
const co = require('co')

describe('akv', function () {
  this.timeout(3000)

  before(() => co(function * () {

  }))

  after(() => co(function * () {

  }))

  it('Akv', () => co(function * () {
    let filename = `${__dirname}/../tmp/testing-akv/akv.json`
    let akv = new AKV(filename, { interval: 100 })

    yield akv.touch()
    assert.ok(fs.existsSync(filename))
    yield akv.destroy()
    assert.ok(!fs.existsSync(filename))

    {
      let { storage } = akv
      yield storage.write({ foo: 'baz' })
      let data = yield storage.read()
      assert.deepEqual(data, { foo: 'baz' })
      assert.ok(storage.needsFlush)
      yield asleep(300)
      assert.ok(!storage.needsFlush)
    }

    for (let i = 0; i < 5; i++) {
      yield akv.set('index', String(i))
      let index = yield akv.get('index')
      assert.equal(index, String(i))
    }

    for (let i = 0; i < 100; i++) {
      akv.set('index', String(i))
      akv.get('index').then((index) => {
        // console.log('index', index)
      })
    }
    assert.equal(yield akv.get('index'), '99')

    assert.deepEqual(yield akv.all(), yield akv.all())
    assert.deepEqual(yield akv.keys(), yield akv.keys())

    yield akv.commit()
  }))
})

/* global describe, before, after, it */
