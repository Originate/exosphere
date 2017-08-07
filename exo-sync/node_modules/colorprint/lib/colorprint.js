/**
 * Colorpint context.
 * @memberof module:colorprint/lib
 * @inner
 * @constructor Colorprint
 * @param {object} config - Context configuration.
 */

'use strict'

const formatMsg = require('./msg/format_msg')
const decorateMsg = require('./msg/decorate_msg')
const indentMsg = require('./msg/indent_msg')

/** @lends module:colorprint/lib~Colorprint */
function Colorprint () {
  const s = this
  s.init.apply(s, arguments)
}

Colorprint.prototype = {
  disabled: false,
  prepareMsg () {
    const s = this
    let msg = formatMsg.apply(formatMsg, arguments)
    return [s.PREFIX, indentMsg(msg, s.indent), s.SUFFIX].join('')
  },
  writeToStdout (msg, color) {
    const s = this
    if (s.disabled) {
      return
    }
    console.log(decorateMsg(msg, color))
  },
  writeToStderr (msg, color) {
    console.error(decorateMsg(msg, color))
  },
  /** Color for point print. */
  POINT_COLOR: 'blue',
  /** Color for notice print. */
  NOTICE_COLOR: 'magenta',
  /** Color for info print. */
  INFO_COLOR: 'green',
  /** Color for debug print. */
  DEBUG_COLOR: '',
  /** Color for trace print. */
  TRACE_COLOR: 'white',
  /** Color for warn print. */
  WARN_COLOR: 'yellow',
  /** Color for error print. */
  ERROR_COLOR: 'red',
  /** Color for fatal print. */
  FATAL_COLOR: 'bgRed',
  /** Alias for module:colorprint/lib~Colorprint#point. */
  POINT () {
    const s = this
    s.point.apply(s, arguments)
  },
  /** Alias for module:colorprint/lib~Colorprint#notice. */
  NOTICE () {
    const s = this
    s.notice.apply(s, arguments)
  },
  /** Alias for module:colorprint/lib~Colorprint#info. */
  INFO () {
    const s = this
    s.info.apply(s, arguments)
  },
  /** Alias for module:colorprint/lib~Colorprint#debug. */
  DEBUG () {
    const s = this
    s.debug.apply(s, arguments)
  },
  /** Alias for module:colorprint/lib~Colorprint#trace. */
  TRACE () {
    const s = this
    s.trace.apply(s, arguments)
  },
  /** Alias for module:colorprint/lib~Colorprint#warn. */
  WARN () {
    const s = this
    s.warn.apply(s, arguments)
  },
  /** Alias for module:colorprint/lib~Colorprint#error. */
  ERROR () {
    const s = this
    s.error.apply(s, arguments)
  },
  /** Alias for module:colorprint/lib~Colorprint#fatal. */
  FATAL () {
    const s = this
    s.fatal.apply(s, arguments)
  },
  /** @constructs module:colorprint/lib~Colorprint */
  init (config = {}) {
    const s = this
    Object.assign(s, config)
    s.PREFIX = config.prefix || s.PREFIX
    s.SUFFIX = config.suffix || s.SUFFIX
  },
  /** Number of indent */
  indent: 0,
  /** Message prefix */
  PREFIX: '',
  /** Message suffix */
  SUFFIX: '',
  /**
   * Print point message.
   * @param {...string} msg - Message to print.
   */
  point (msg) {
    const s = this
    s.writeToStdout(s.prepareMsg.apply(s, arguments), s.POINT_COLOR)
  },
  /**
   * Print notice message.
   * @param {...string} msg - Message to print.
   */
  notice (msg) {
    const s = this
    s.writeToStdout(s.prepareMsg.apply(s, arguments), s.NOTICE_COLOR)
  },
  /**
   * Print info message.
   * @param {...string} msg - Message to print.
   */
  info (msg) {
    const s = this
    s.writeToStdout(s.prepareMsg.apply(s, arguments), s.INFO_COLOR)
  },
  /**
   * Print debug message.
   * @param {...string} msg - Message to print.
   */
  debug (msg) {
    const s = this
    s.writeToStdout(s.prepareMsg.apply(s, arguments), s.DEBUG_COLOR)
  },
  /**
   * Print trace message.
   * @param {...string} msg - Message to print.
   */
  trace (msg) {
    const s = this
    s.writeToStdout(s.prepareMsg.apply(s, arguments), s.TRACE_COLOR)
  },
  /**
   * Print warn message.
   * @param {...string} msg - Message to print.
   */
  warn (msg) {
    const s = this
    s.writeToStdout(s.prepareMsg.apply(s, arguments), s.WARN_COLOR)
  },
  /**
   * Print error message.
   * @param {...string} msg - Message to print.
   */
  error (msg) {
    const s = this
    s.writeToStderr(s.prepareMsg.apply(s, arguments), s.ERROR_COLOR)
  },
  /**
   * Print fatal message.
   * @param {...string} msg - Message to print.
   */
  fatal (msg) {
    const s = this
    s.writeToStderr(s.prepareMsg.apply(s, arguments), s.FATAL_COLOR)
  }
}

module.exports = Colorprint
