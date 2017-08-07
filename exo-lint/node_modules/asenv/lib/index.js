/**
 * NODE_ENV accessor
 * @module asenv
 */

'use strict'

let d = (module) => module.default || module

module.exports = {
  get constants () { return d(require('./constants')) },
  get getEnv () { return d(require('./get_env')) },
  get isDevelopment () { return d(require('./is_development')) },
  get isProduction () { return d(require('./is_production')) },
  get isTest () { return d(require('./is_test')) },
  get setEnv () { return d(require('./set_env')) }
}
