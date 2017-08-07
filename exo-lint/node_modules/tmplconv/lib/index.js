/**
 * Two way template converter.
 * @module tmplconv
 */

'use strict'

let d = (module) => module && module.default || module

module.exports = {
  get render () { return d(require('./render')) },
  get tmplify () { return d(require('./tmplify')) },
  get transplant () { return d(require('./transplant')) }
}
