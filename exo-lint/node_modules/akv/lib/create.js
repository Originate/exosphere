/**
 * Create an AKV instance
 * @function create
 */
'use strict'

const AKV = require('./akv')

/** @lends create */
function create (...args) {
  return new AKV(...args)
}

module.exports = create
