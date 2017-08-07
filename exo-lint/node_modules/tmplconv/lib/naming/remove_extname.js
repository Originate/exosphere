/**
 * Remove file extension.
 * @function removeExtname
 * @param {string} filename - Filename to change.
 * @param {string} extname - Extname to remove.
 * @returns {string} Changed file name.
 */

'use strict'

const path = require('path')

/** @lends removeExtname */
function removeExtname (filename, extname) {
  let dirname = path.dirname(filename)
  let basename = path.basename(filename, extname)
  return path.join(dirname, basename)
}

module.exports = removeExtname;
