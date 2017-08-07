/**
 * Generate template from rendered.
 * @memberof module:tmplconv/lib
 * @function tmplify
 * @param {string} srcDir - Name of destination directory.
 * @param {string} destDir - Name of destination directory, which contains template files.
 * @param {object} options - Optional settings.
 * @param {string|object} options.data - Data for tmpls.
 * @param {string|string[]} options.pattern - Source patterns.
 * @param {string|string[]} [options.ignore] - Filename pattern.
 * @param {string} [options.prefix='_____'] - Embed prefix.
 * @param {string} [options.suffix='_____'] - Embed suffix.
 * @param {string} [options.extname='.tmpl'] - Template extension name.
 * @param {boolean} [options.silent=false] - Silent or not.
 * @param {string} [options.mode='644'] - File permission to generate.
 * @param {boolean} [options.clean=false] - Cleanup destination directory before convert.
 * @param {boolean} [options.once=false] - Write only first time. Skip if already exists.
 * @returns {Promise}
 */

'use strict'

const argx = require('argx')
const convert = require('./converting/convert')
const addExtname = require('./naming/add_extname')
const _rule = require('./_rule')
const _logResults = require('./_log_results')
const co = require('co')

/** @lends tmplify */
function tmplify (srcDir, destDir, options) {
  let args = argx(arguments)
  if (args.pop('function')) {
    throw new Error('Callback is no longer supported. Use promise interface instead.')
  }
  options = args.pop('object') || {}

  return co(function * () {
    let results = yield convert(srcDir, destDir, {
      pattern: options.pattern || '**/*.*',
      rule: _rule.reversed(options.data, options.prefix, options.suffix),
      ignore: options.ignore || 'node_modules/**/*.*',
      mode: options.mode,
      clean: options.clean,
      once: options.once,
      out (src) {
        return addExtname(src, options.extname || '.tmpl')
      }
    })
    let shouldLog = !options.silent
    if (shouldLog) {
      _logResults(results)
    }

    // Wait for flush
    yield new Promise((resolve) =>
      process.nextTick(() => resolve())
    )

    return results
  })
}

module.exports = tmplify
