/**
 * Check if the env is production
 * @function isProduction
 * @returns {boolean} - Production or not
 */
'use strict'

const { PRODUCTION } = require('./constants')
const getEnv = require('./get_env')

/** @lends isProduction */
function isProduction () {
  return getEnv() === PRODUCTION
}

module.exports = isProduction
