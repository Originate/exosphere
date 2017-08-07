/**
 * Create a colorprint instance.
 * @memberof module:colorprint/lib
 * @function create
 * @param {object} config - Colorprint config.
 * @returns {Colorprint} - Colorprint instance.
 */

'use strict'

const Colorprint = require('./colorprint')

/** @lends create */
function create (config) {
  return new Colorprint(config)
}

module.exports = create
