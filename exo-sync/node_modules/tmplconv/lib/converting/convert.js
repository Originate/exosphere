/**
 * Convert files in a directory.
 * @function convertDir
 * @param {string} src - Source directory path.
 * @param {string} dest - Destination directory path.
 * @param {object} [options] - Optional settings.
 * @param {string|string[]} [options.pattern] - Filename pattern.
 * @param {string|string[]} [options.ignore] - Filename pattern.
 * @param {string} [options.mode='644'] - File permission to generate.
 * @param {boolean} [options.once=false] - Write only first time. Skip if already exists.
 * @param {object} [options.rule] - Convert map.
 * @returns {Promise}
 */

'use strict'

const co = require('co')
const { existsAsync, statAsync } = require('asfs')
const convertDir = require('./convert_dir')
const convertFile = require('./convert_file')

/** @lends convert */
function convert (src, dest, options = {}) {
  return co(function * () {
    let exists = yield existsAsync(src)
    if (!exists) {
      throw new Error(`src not exists: ${src}`)
    }
    let stats = yield statAsync(src)
    let isDir = stats.isDirectory()
    if (isDir) {
      return yield convertDir(src, dest, options)
    } else {
      return yield convertFile(src, dest, options)
    }
  })
}

module.exports = convert
