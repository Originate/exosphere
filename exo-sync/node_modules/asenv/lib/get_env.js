/**
 * Get NODE_ENV value
 * @function getEnv()
 * @returns {?string} - Env value
 */
'use strict'

/** @lends getEnv */
function getEnv () {
  return process.env.NODE_ENV
}

module.exports = getEnv
