/**
 * @augments AKV
 * @class AKVStatus
 */
'use strict'

const { AKV, fileHash } = require('akv')
const aglob = require('aglob')
const co = require('co')

/** @lends AKVStatus */
class AKVStatus extends AKV {
  constructor (filename, options) {
    super(filename, options)
  }

  /**
   * Detect unknown status data
   * @param {string[]} filenames
   * @returns {Promise}
   */
  filterStatusUnknown (filenames) {
    const s = this
    return co(function * () {
      let filtered = []
      for (let filename of filenames) {
        let knownHash = yield s.get(filename)
        let hash = yield fileHash(filename)
        let known = hash && knownHash && (hash === knownHash)
        if (!known) {
          filtered.push(filename)
        }
      }
      return filtered
    })
  }

  /**
   * @param {string[]} filenames
   * @returns {*|Promise}
   */
  saveStatus (filenames) {
    const s = this
    return co(function * () {
      for (let filename of filenames) {
        let hash = yield fileHash(filename)
        yield s.set(filename, hash)
      }
      yield s.commit()
    })
  }
}

module.exports = AKVStatus
