/**
 * A simple key value store using single json file
 * @module akv
 */

'use strict'

let d = (module) => module.default || module

module.exports = {
  get fileHash () { return d(require('./file_hash')) },
  get readFromFile () { return d(require('./read_from_file')) },
  get writeToFile () { return d(require('./write_to_file')) }
}
