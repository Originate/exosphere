/**
 * Convert string.
 * @function convertString
 * @param {string} source - Source string to convert.
 * @param {object} rule - Convert map.
 * @returns {string} - Converted string.
 */

'use strict'

/** @lends convertString */
function convertString (source, rule) {
  let result = String(source)
  Object.keys(rule).forEach((key) => {
    let regExp = new RegExp(key, 'g')
    let replacings = [].concat(rule[ key ])
    replacings.map((replacing) => {
      result = result.replace(regExp, replacing)
    })
  })
  return result
}

module.exports = convertString
