/**
 * Register babel
 * @function registerES2015
 * @param {Object} options
 */
'use strict'

const register = require('ababel/register')

/** @lends registerES2015 */
function registerES2015 (options) {
  options = options || {}
  options.presets = [ 'es2015' ].concat(options.presets || [])
  register(options)
}

module.exports = registerES2015
