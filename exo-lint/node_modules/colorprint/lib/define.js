/**
 * Define a logger.
 * @memberof module:colorprint/lib
 * @function define
 * @param {object} properties - Logger prototype properties.
 * @returns {function} - A logger constructor.
 */

'use strict'

const create = require('./create')

/** @lends define */
function define (properties) {
  function Logger (config = {}) {
    const s = this
    Object.assign(s, config)
    s.PREFIX = config.prefix || s.PREFIX
    s.SUFFIX = config.suffix || s.SUFFIX
  }

  Logger.prototype = create(properties)
  return Logger
}

module.exports = define
