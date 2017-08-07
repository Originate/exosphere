/**
 * Decorate message with ansi color.
 * @see https://www.npmjs.com/package/cli-color
 * @memberof module:colorprint/lib
 * @function decorateMsg
 * @param {string} msg - Messages to decorateMsg.
 * @param {string} color - Name of color.
 * @returns {string}  - Decorated message.
 */

'use strict'

const cliColor = require('cli-color')

/** @lends decorateMsg */
function decorateMsg (msg, color) {
  if (!color) {
    return msg
  }
  let decorator = color && cliColor[ color ]
  if (!decorator) {
    throw new Error('Unknown color: ' + color)
  }
  return decorator(msg)
}

module.exports = decorateMsg

