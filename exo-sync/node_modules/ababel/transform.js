/**
 * Browserify transform
 * @function transform
 * @param {Object} options
 */
'use strict'

const babelify = require('babelify')
const constants = require('./lib/constants')

/** @lends transform */
function transform (options) {
  return babelify.configure(
    Object.assign({
      extensions: constants.DEFAULT_EXT.split(','),
      compact: false,
      babelrc: false,
      sourceRoot: process.cwd(),
      presets: constants.DEFAULT_PRESET.split(',')
    }, options || {})
  )
}

module.exports = transform
