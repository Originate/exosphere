/**
 * Test case for colorprint bin.
 * Runs with mocha.
 */
'use strict'

const bin = require.resolve('../bin/colorprint')
const assert = require('assert')
const childProcess = require('child_process')

describe('bin', ()=> {
  function _spawn (command, args) {
    let spawned = childProcess.spawn(command, args)
    spawned.stdout.pipe(process.stdout)
    spawned.stderr.pipe(process.stderr)
  }

  it('Print.', (done) => {
    _spawn(bin, [ 'notice', 'This is notice', 'from cli.' ])
    _spawn(bin, [ 'info', 'This is info', 'from cli.' ])
    _spawn(bin, [ 'debug', 'This is debug', 'from cli.' ])
    _spawn(bin, [ 'trace', 'This is trace', 'from cli.' ])
    _spawn(bin, [ 'error', 'This is error', 'from cli.' ])
    _spawn(bin, [ 'fatal', 'This is fatal', 'from cli.' ])
    done()
  })
})

/* global describe, it */
