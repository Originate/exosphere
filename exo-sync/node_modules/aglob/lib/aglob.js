/**
 * Async glob
 * @function aglob
 * @param {string|string[]} patterns - Pattern
 * @param {Object} [options] - Optional settings
 * @returns {Promise}
 */
'use strict'

const co = require('co')
const glob = require('glob')

/** @lends aglob */
function aglob (patterns, options = {}) {
  return co(function * () {
    let results = []
    for (let pattern of [].concat(patterns || [])) {
      let filenames = yield new Promise((resolve, reject) =>
        glob(pattern, options, (err, filenames) =>
          err ? reject(err) : resolve(filenames)
        )
      )
      results = results.concat(filenames)
    }
    return results
  })
}

Object.assign(aglob, {
  sync (patterns, options) {
    return [].concat(patterns || [])
      .reduce((result, pattern) => result.concat(glob.sync(pattern, options)), [])
  }
})

module.exports = aglob
