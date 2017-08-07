/**
 * @function _hasDuplicate
 * @returns {Promsie}
 */

'use strict'

const fs = require('fs')
const co = require('co')

/** @lends _hasDuplicate */
function _hasDuplicate (filename, content) {
  return co(function * () {
    let exists = yield new Promise((resolve) =>
      fs.exists(filename, (exists) => resolve(exists))
    )
    if (!exists) {
      return false
    }
    let existing = yield new Promise((resolve, reject) =>
      fs.readFile(filename, (err, content) =>
        err ? reject(err) : resolve(content)
      )
    )
    return String(existing) === String(content)
  })
}

module.exports = _hasDuplicate
