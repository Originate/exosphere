/**
 * Compile files
 * @function ababel
 * @param {string} pattern - Glob file name pattern
 * @param {Object} [options] - Optional settings
 * @param {string} [options.status] - Status file path
 * @param {string} [options.cwd] - Current working directory path
 * @param {string} [options.out] - Output directory path
 * @param {boolean} [options.minified] - Minified or not
 * @param {string[]} [options.reflects] - File patterns to reflects changes
 * @returns {Promise}
 */
'use strict'

const aglob = require('aglob')
const convertSourceMap = require('convert-source-map')
const writeout = require('writeout')
const akvStatus = require('akv-status')
const co = require('co')
const path = require('path')
const filedel = require('filedel')
const { statAsync } = require('asfs')
const { isProduction } = require('asenv')
const { transformFile } = require('babel-core')
const { DEFAULT_PRESET, DEFAULT_EXT } = require('./constants')

let { replaceExt } = require('./helpers/file_helper')
let relative = (filename) => path.relative(process.cwd(), filename)
let mtime = (filename) => statAsync(filename).catch(() => null).then(({mtime}) => mtime)

/** @lends ababel */
function ababel (pattern, options = {}) {
  let {
    status = 'tmp/ababel.status.json',
    cwd = process.cwd(),
    out = process.cwd(),
    presets = DEFAULT_PRESET.split(','),
    sourceMaps = !isProduction(),
    minified = false,
    ignore = [],
    reflects = [],
    ext = DEFAULT_EXT.split(',')
  } = options

  let store = akvStatus(status)
  return co(function * () {
    let filenames = yield aglob(pattern, { cwd, ignore })
    reflects = yield aglob(reflects)
    for (let filename of filenames) {
      let src = path.resolve(cwd, filename)
      let dest = path.resolve(out, ext.reduce((filename, ext) => replaceExt(filename, ext, '.js'), filename))
      if (!isProduction()) {
        let changed = yield store.filterStatusUnknown([ src, dest, ...reflects ])
        if (changed.length === 0) {
          let srcMtime = yield mtime(src)
          let destMtime = yield mtime(dest)
          let skip = srcMtime && destMtime && (srcMtime <= destMtime)
          if (skip) {
            continue
          }
        }
      }
      let { code, map, ast } = yield new Promise((resolve, reject) => {
        let options = {
          presets,
          minified,
          sourceMaps,
          compact: false,
          babelrc: false,
          sourceRoot: path.relative(process.cwd(), cwd),
          plugins: [ 'transform-runtime' ]
        }
        transformFile(src, options, (err, result) => err ? reject(err) : resolve(result))
      })
      try {
        if (dest !== src) {
          yield filedel(dest)
        }
      } catch (err) {
        // Do nothing
      }
      let { skipped } = yield writeout(dest, `${code}\n${convertSourceMap.fromObject(map).toComment()}`, {
        mkdirp: true,
        skipIfIdentical: true
      })
      if (!skipped) {
        console.log(`File generated: ${relative(dest)}`)
      }

      yield store.saveStatus([ src, dest, ...reflects ])
    }
  }).catch((err) => co(function * () {
    yield store.destroy()
    throw err
  }))
}

module.exports = ababel
