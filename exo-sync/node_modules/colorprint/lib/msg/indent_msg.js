/**
 * Add indent.
 * @function indentMsg
 * @param {string} msg - Message format.
 * @param {number} depth - Depth to indent.
 * @returns {string} - Indented message.
 */

'use strict'

const os = require('os')

const { EOL } = os
const TAB = '  '

/** @lends indentMsg */
function indentMsg (msg, depth) {
  let indent = ''
  for (let i = 0; i < depth; i++) {
    indent += TAB
  }
  return indent + msg.replace(new RegExp(EOL, 'g'), EOL + indent)
}

module.exports = indentMsg
