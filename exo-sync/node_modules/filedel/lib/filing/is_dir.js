/**
 * @function isDir
 * @returns {Promise}
 */

'use strict'

const { existsAsync, statAsync } = require('asfs')
const co = require('co')

/** @lends isDir */
function isDir (filename) {
  return co(function * () {
    let exists = yield existsAsync(filename)
    if (!exists) {
      return
    }
    try {
      let stats = yield statAsync(filename)
      return stats.isDirectory()
    } catch (e) {
      // Ignore error.
    }
  })
}

module.exports = isDir
