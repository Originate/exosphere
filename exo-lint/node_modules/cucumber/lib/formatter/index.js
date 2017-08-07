'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var Formatter = function Formatter(options) {
  (0, _classCallCheck3.default)(this, Formatter);

  _lodash2.default.assign(this, _lodash2.default.pick(options, ['colorFns', 'cwd', 'log', 'snippetBuilder', 'stream', 'supportCodeLibrary']));
};

exports.default = Formatter;