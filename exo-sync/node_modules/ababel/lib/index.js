/**
 * Compile with babel only when file changed from last time.
 * @module ababel
 */

'use strict'

const ababel = require('./ababel')

let lib = ababel.bind(this)

Object.assign(lib, ababel, {
  ababel
})

module.exports = lib