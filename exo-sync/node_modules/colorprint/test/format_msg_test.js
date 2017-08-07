/**
 * Test case for formatMsg.
 * Runs with mocha.
 */
'use strict'

const formatMsg = require('../lib/msg/format_msg.js')
const assert = require('assert')

describe('format', () => {
  it('Format msg', (done) => {
    assert.equal(formatMsg('Hey, my name is %s, I am %d years old.', 'John', 34, 'Hoo!'), "Hey, my name is John, I am 34 years old. Hoo!")
    assert.equal(formatMsg(''), '')
    assert.equal(formatMsg(), '')
    assert.equal(formatMsg('foo%f', 0.4), 'foo0.4')
    assert.equal(formatMsg('foo%j', 0.4), 'foo%j 0.4')
    done()
  })

  it('Format msg with object', (done) => {
    let msg = formatMsg({ foo: 'bar' }, null)
    assert.ok(msg)

    done()
  })
})

/* global describe, it */
