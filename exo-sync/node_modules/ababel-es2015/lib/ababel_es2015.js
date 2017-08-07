/**
 * Compile es2015 javascripts
 * @function ababelES2015
 * @param {string} pattern - Glob file name pattern
 * @param {Object} [options] - Optional settings
 */
'use strict'

const ababel = require('ababel')

/** @lends ababelES2015 */
function ababelES2015 (pattern, options = {}) {
  let { presets } = options
  options.presets = [ 'es2015' ].concat(presets || [])
    .filter((name, i, array) => array.indexOf(name) === i)
  return ababel(pattern, options)
}

module.exports = ababelES2015
