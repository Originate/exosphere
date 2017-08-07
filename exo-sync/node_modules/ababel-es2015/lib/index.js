/**
 * ababel with es2015 flavor
 * @module ababel-es2015
 * @version 1.0.8
 */

'use strict'

const ababelES2015 = require('./ababel_es2015')

let lib = ababelES2015.bind(this)

Object.assign(lib, ababelES2015, {
  ababelES2015
})

module.exports = lib