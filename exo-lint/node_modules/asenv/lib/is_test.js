/**
 * Check if the env is test
 * @function isTest
 * @returns {boolean} - Test or not
 */
'use strict'

const { TEST } = require('./constants')
const getEnv = require('./get_env')

/** @lends isTest */
function isTest () {
  return getEnv() === TEST
}

module.exports = isTest
