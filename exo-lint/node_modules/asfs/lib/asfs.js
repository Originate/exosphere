/**
 * @module asfs
 */
'use strict'

const fs = require('fs')
const mkdirp = require('mkdirp')

module.exports = Object.assign(exports, {
  /**
   * Check file exists
   * @param {string} filename
   * @returns {Promise.<boolean>}
   */
  existsAsync (filename) {
    return new Promise((resolve) =>
      fs.exists(filename, (exists) => resolve(exists))
    )
  },

  /**
   * Make directories with parent option.
   * @param {string} dirname
   * @returns {Promise}
   */
  mkdirpAsync (dirname) {
    return new Promise((resolve, reject) =>
      mkdirp(dirname, (err) => err ? reject(err) : resolve())
    )
  },

  /**
   * Read file
   * @param {string} filename
   * @param {string} [encode]
   * @returns {Promise}
   */
  readFileAsync (filename, encode) {
    let args = [ ...arguments ]
    return new Promise((resolve, reject) =>
      fs.readFile(...args.concat((err, content) => err ? reject(err) : resolve(content)))
    )
  },

  /**
   * Write a file
   * @param {string} filename
   * @param {string} content
   * @returns {Promise}
   */
  writeFileAsync (filename, content) {
    let args = [ ...arguments ]
    return new Promise((resolve, reject) =>
      fs.writeFile(...args.concat((err) => err ? reject(err) : resolve()))
    )
  },

  /**
   * Unlink a file
   * @param {string} filename
   * @returns {Promise}
   */
  unlinkAsync (filename) {
    return new Promise((resolve, reject) =>
      fs.unlink(filename, (err) => err ? reject(err) : resolve())
    )
  },

  /**
   * Get state of a file
   * @param {string} filename
   * @returns {Promise}
   */
  statAsync (filename) {
    return new Promise((resolve, reject) =>
      fs.stat(filename, (err, stat) => err ? reject(err) : resolve(stat))
    )
  },

  /**
   * Read directory
   * @param {string} dirname
   * @param {Object} [options]
   * @returns {Promise}
   */
  readdirAsync (dirname, options = {}) {
    return new Promise((resolve, reject) =>
      fs.readdir(dirname, options, (err, result) => err ? reject(err) : resolve(result))
    )
  }
})
