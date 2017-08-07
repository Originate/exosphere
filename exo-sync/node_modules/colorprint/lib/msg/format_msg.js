/**
 * Format message.
 * @memberof module:colorprint/lib
 * @function formatMsg
 * @param {...string} msg - Messages to format.
 * @returns {string}  - Formatted message.
 */

"use strict"

/** @lends formatMsg */
function formatMsg (msg) {
  const s = this
  msg = Array.prototype.slice.call(arguments, 0).map(msg => {
    if (typeof(msg) === 'object') {
      try {
        return JSON.stringify(msg, null, 2)
      } catch (e) {
        // Do nothing.
      }
    }
    return msg
  })
    .filter(msg => !!msg)
    .map(msg => String(msg))
  if (!msg.length) {
    return ''
  }
  let formatted = msg.shift().replace(/%(.)/g, ($0, $1) => {
    switch ($1) {
      case 's':
        return String(msg.shift())
      case 'd':
        return parseInt(msg.shift())
      case 'f':
        return parseFloat(msg.shift())
      default:
        return $0
    }
  })
  return [ formatted ].concat(msg).join(' ')
}

module.exports = formatMsg
