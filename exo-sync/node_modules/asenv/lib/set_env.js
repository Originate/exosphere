/**
 * Set NODE_ENV value
 * @function setEnv()
 * @param {string} env - Env value
 */
'use strict'

/** @lends setEnv */
function setEnv (env) {
  Object.assign(process.env, {
    NODE_ENV: env && String(env).trim() || env
  })
}

module.exports = setEnv
