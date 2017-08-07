/**
 * Create an AKVStatus instance
 * @function create
 */
'use strict'

const AKVStatus = require('./akv_status')

/** @lends create */
function create (...args) {
  return new AKVStatus(...args)
}

module.exports = create
