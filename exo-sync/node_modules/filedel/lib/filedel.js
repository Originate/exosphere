/**
 * Delete file.
 * @function filedel
 * @param {string} filename - Filename to delete.
 * @param {object} [options] - Optional settings.
 * @param {boolean} [options.force=false] - Unlink even if readonly.
 * @returns {Promise}
 */

'use strict'

const co = require('co')
const rimraf = require('rimraf')
const { existsAsync } = require('asfs')
const argx = require('argx')
const aglob = require('aglob')
const doUnlink = require('./filing/do_unlink')
const isDir = require('./filing/is_dir')

/** @lends filedel */
function filedel (target, options = {}) {
  let args = argx(arguments)
  if (args.pop('function')) {
    throw new Error('[filedel] Callback is no more supported. Use promise interface instead.')
  }
  options = args.pop('object') || {}

  return co(function * () {
    let filenames = yield aglob(target, {})
    for (let filename of filenames) {
      let exists = yield existsAsync(filename)
      if (!exists) {
        return
      }
      let isDir_ = yield isDir(filename)
      if (isDir_) {
        throw new Error(`[filedel] Can not unlink directory: ${filename}`)
      }
      yield doUnlink(filename, !!options.force)
    }
  })
}

Object.assign(filedel, {
  recursive (dirname) {
    return co(function * () {
      let exists = yield existsAsync(dirname)
      if (!exists) {
        return
      }
      let isDir_ = yield isDir(dirname)
      if (!isDir_) {
        throw new Error(`[filedel] Not a directory: ${dirname}`)
      }
      yield new Promise((resolve, reject) =>
        rimraf(dirname, (err) => err ? reject(err) : resolve())
      )
    })
  }
})

module.exports = filedel
