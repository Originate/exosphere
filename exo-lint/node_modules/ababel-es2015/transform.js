/**
 * Browserify transform
 * @function transformES2015
 * @param {Object} options
 */
'use strict'

const transform = require('ababel/transform')

/** @lends transformES2015 */
function transformES2015 (options) {
  options = options || {}
  options.presets = [ 'es2015' ].concat(options.presets || [])
  return transform(options)
}

module.exports = transformES2015
