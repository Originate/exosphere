/**
 * Log results.
 * @function _logResults
 * @private
 */

'use strict'

const colorprint = require('colorprint')

/** @lends _logResults */
function _logResults (results) {
  [].concat(results).forEach((result) => {
    let hasResult = result && !result.skipped
    if (hasResult) {
      colorprint.debug('File generated: ' + result.filename)
    }
  })
}

module.exports = _logResults
