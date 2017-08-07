/**
 * @function doUnlink
 * @returns {Promise}
 */
'use strict'

const fs = require('fs')
const co = require('co')

/** @lends doUnlink */
function doUnlink (filename, force) {
  return co(function * () {
    if (force) {
      yield new Promise((resolve, reject) =>
        fs.chmod(filename, '666', (err) =>
          err ? reject(err) : resolve()
        )
      )
    }
    yield new Promise((resolve, reject) =>
      fs.unlink(filename, (err) =>
        err ? reject(err) : resolve()
      )
    )
  })
}

module.exports = doUnlink
