/**
 * Writeout a file.
 * @function writeout
 * @param {string} filename - Filename to write.
 * @param {string} content - String content to write.
 * @param {object} [options] - Optional settings.
 * @param {boolean} [options.mkdirp=false] - Make parent directories.
 * @param {boolean} [options.skipIfIdentical=false] - Skip to write if existing content is identical.
 * @param {string} [options.mode='644'] - File permission.
 * @param {boolean} [options.force=false] - Force to write even if existing file is readonly.
 * @returns {Promise}
 */

'use strict'

const co = require('co')
const fs = require('fs')
const path = require('path')
const argx = require('argx')
const mkdirp = require('mkdirp')
const filedel = require('filedel')
const _hasDuplicate = require('./_has_duplicate')

/** @lends writeout */
function writeout (filename, content, options = {}) {
  if (argx(arguments).pop('function')) {
    throw new Error('[writeout] Callback is no more supported. Use promise interface instead.')
  }

  let needsMkdirp = !!options.mkdirp
  let skipIfIdentical = !!options.skipIfIdentical
  let force = !!options.force
  let mode = options.mode || '644'

  let result = { filename }

  return co(function * () {
    if (needsMkdirp) {
      yield new Promise((resolve, reject) =>
        mkdirp(path.dirname(filename), (err) => err ? reject(err) : resolve())
      )
    }

    if (skipIfIdentical) {
      let skip = yield _hasDuplicate(filename, content)
      if (skip) {
        result.skipped = true
        return result
      }
    }

    if (force) {
      yield filedel(filename, { force: true })
    }

    yield new Promise((resolve, reject) =>
      fs.writeFile(filename, content, {
        mode: mode
      }, (err) => err ? reject(err) : resolve())
    )

    result.skipped = false

    // Wait for flush
    yield new Promise((resolve) =>
      process.nextTick(() => resolve())
    )
    
    return result
  })
}

module.exports = writeout
