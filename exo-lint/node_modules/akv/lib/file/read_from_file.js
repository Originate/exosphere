/**
 * Read a json file
 * @function readFromFile
 * @param {string} filename - File name to read
 * @returns {Promise.<Object>} - Data json file
 */
'use strict'

const co = require('co')
const { existsAsync, readFileAsync } = require('asfs')

/** @lends readFromFile */
function readFromFile (filename) {
  return co(function * () {
    let exists = yield existsAsync(filename)
    if (!exists) {
      return null
    }
    let content = (yield readFileAsync(filename)).toString()
    if (!content) {
      return null
    }
    return JSON.parse(content)
  })
}

module.exports = readFromFile
