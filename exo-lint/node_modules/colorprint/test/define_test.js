/**
 * Test case for define.
 * Runs with mocha.
 */
'use strict'

const define = require('../lib/define.js')
const assert = require('assert')

describe('define', ()=> {
  it('Define', (done) => {
    let Logger = define({
      verbose: false
    })
    let logger = new Logger({})
    logger.INFO('This is custom INFO')
    logger.DEBUG('This is custom DEBUG')
    logger.TRACE('This is custom TRACE')
    logger.ERROR('This is custom ERROR')
    logger.WARN('This is custom WARN')
    logger.FATAL('This is custom FATAL')
    done()
  })
})

/* global describe, it */
