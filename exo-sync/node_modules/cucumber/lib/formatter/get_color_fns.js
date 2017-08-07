'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _defineProperty2 = require('babel-runtime/helpers/defineProperty');

var _defineProperty3 = _interopRequireDefault(_defineProperty2);

exports.default = getColorFns;

var _safe = require('colors/safe');

var _safe2 = _interopRequireDefault(_safe);

var _status = require('../status');

var _status2 = _interopRequireDefault(_status);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function getColorFns(enabled) {
  var _colors$setTheme;

  _safe2.default.enabled = enabled;
  _safe2.default.setTheme((_colors$setTheme = {}, (0, _defineProperty3.default)(_colors$setTheme, _status2.default.AMBIGUOUS, 'red'), (0, _defineProperty3.default)(_colors$setTheme, _status2.default.FAILED, 'red'), (0, _defineProperty3.default)(_colors$setTheme, _status2.default.PASSED, 'green'), (0, _defineProperty3.default)(_colors$setTheme, _status2.default.PENDING, 'yellow'), (0, _defineProperty3.default)(_colors$setTheme, _status2.default.SKIPPED, 'cyan'), (0, _defineProperty3.default)(_colors$setTheme, _status2.default.UNDEFINED, 'yellow'), (0, _defineProperty3.default)(_colors$setTheme, 'location', 'grey'), (0, _defineProperty3.default)(_colors$setTheme, 'tag', 'cyan'), _colors$setTheme));
  return _safe2.default;
}