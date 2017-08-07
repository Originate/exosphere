#!/usr/bin/env node

/**
 * Prepare shims
 */

'use strict'

process.chdir(`${__dirname}/..`)

const apeTasking = require('ape-tasking')
const ababel = require('../lib')

apeTasking.runTasks('shim', [
  () => ababel('**/*.js', {
    cwd: 'lib',
    out: 'shim/node',
    presets: [ 'es2015' ]
  })
], true)
