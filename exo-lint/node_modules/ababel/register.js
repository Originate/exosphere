/**
 * Register babel
 * @function register
 * @param {Object} options
 */
'use strict'

const babelRegister = require('babel-register')
const constants = require('./lib/constants')

/** @lends register */
function register (options) {
  babelRegister(
    Object.assign({
      extensions: constants.DEFAULT_EXT.split(','),
      compact: false,
      babelrc: false,
      sourceRoot: process.cwd(),
      presets: constants.DEFAULT_PRESET.split(',')
    }, options || {})
  )
}

module.exports = register
