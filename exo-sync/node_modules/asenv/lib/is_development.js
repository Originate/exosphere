/**
 * Check if the env it development
 * @function isDevelopment
 * @returns {boolean} - Development or not
 */
'use strict'

const { DEVELOPMENT } = require('./constants')
const getEnv = require('./get_env')

/** @lends isDevelopment */
function isDevelopment () {
  return getEnv() === DEVELOPMENT
}

module.exports = isDevelopment
