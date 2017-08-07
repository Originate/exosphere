/**
 * A simple key value store
 * @class AKV
 * @param {string} filename - Filename to store data
 * @param {Object} [options]
 * @param {number} [options.interval=1000] - Flush interval
 */
'use strict'

const Storage = require('./storage')
const co = require('co')

/** @lends AKV */
class AKV {
  constructor (filename, options = {}) {
    const s = this
    if (!filename) {
      throw new Error('filename is required')
    }
    let {
      interval = 1000
    } = options
    s.storage = new Storage(filename, { interval }).start(interval)

    process.setMaxListeners(process.getMaxListeners() + 2)
    process.on('beforeExit', () => s.handleBeforeExit())
    process.on('exit', () => s.handleExit())
  }

  /**
   * Touch file
   * @function touch
   * @returns {Promise}
   */
  touch () {
    const s = this
    let { storage } = s
    return co(function * () {
      let data = yield storage.read()
      data = data || {}
      yield storage.write(data)
      yield storage.flush()
    })
  }

  /**
   * Set a value
   * @function set
   * @param {string} key
   * @param {string} value
   * @returns {Promise}
   */
  set (key, value) {
    const s = this
    let { storage } = s
    return co(function * () {
      let data = yield storage.read()
      data = data || {}
      data[ key ] = value
      yield storage.write(data)
    })
  }

  /**
   * Get all keys
   * @function keys
   * @returns {Promise}
   */
  keys () {
    const s = this
    let { storage } = s
    return co(function * () {
      let data = yield storage.read()
      return Object.keys(data || {})
    })
  }

  /**
   * Get a value
   * @function get
   * @param {string} key
   * @returns {Promise}
   */
  get (key) {
    const s = this
    let { storage } = s
    return co(function * () {
      let data = yield storage.read()
      return data && data[ key ]
    })
  }

  /**
   * Get all values
   * @function all
   * @returns {Promise}
   */
  all () {
    const s = this
    let { storage } = s
    return co(function * () {
      let data = yield storage.read()
      return data || {}
    })
  }

  /**
   * Delete a value
   * @function del
   * @param {string} key
   * @returns {Promise}
   */
  del (key) {
    const s = this
    let { storage } = s
    return co(function * () {
      let data = yield storage.read()
      if (!data) {
        return
      }
      delete data[ key ]
      yield storage.write(data)
    })
  }

  /**
   * Delete all values
   * @function destroy
   * @returns {Promise}
   */
  destroy () {
    const s = this
    let { storage } = s
    return co(function * () {
      yield storage.purge()
    })
  }

  /**
   * Commit changes
   * @returns {Promise}
   */
  commit () {
    const s = this
    let { storage } = s
    return co(function * () {
      yield storage.flush()
      storage.needsFlush = false
    })
  }

  handleBeforeExit () {
    const s = this
    let { storage } = s
    storage.flushIfNeeded()
    storage.needsFlush = false
  }

  handleExit () {
    const s = this
    let { storage } = s
    if (s.storage.needsFlush) {
      console.warn('[akv] Some uncommitted change has lost. Make sure to call akv.commit() before existing.')
    }
  }
}

module.exports = AKV
