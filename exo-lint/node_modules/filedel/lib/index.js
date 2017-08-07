/**
 * Delete files.
 * @module filedel
 * @version 2.0.5
 */

'use strict'

const filedel = require('./filedel')

let lib = filedel.bind(this)

Object.assign(lib, filedel, {
  filedel
})

module.exports = lib
